package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/samber/lo"
)

//go:embed  system_prompt.md
var systemPrompt string

func main() {
	var ctx = context.Background()
	c, err := NewClient(systemPrompt)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// c.Chat(ctx, "请帮我统计上周券码核销率与这周券码核销率的环比增长率")
	// c.Chat(ctx, "请列出时间处理示例")
	c.Chat(ctx, "请告诉我上周的时间范围")

	time.Sleep(time.Second * 2)
}

type Client struct {
	SystemPrompt string
	MCP          *client.Client
	LLMClient    openai.Client
	ToolList     []openai.ChatCompletionToolParam
}

func NewClient(prompt string) (*Client, error) {
	var ctx = context.Background()
	var httpURL = "http://localhost:8080/mcp"
	httpTransport, err := transport.NewStreamableHTTP(httpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP transport: %v", err)
	}

	c := client.NewClient(httpTransport)
	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	chatClient := &Client{
		SystemPrompt: prompt,
		MCP:          c,
		LLMClient: openai.NewClient(
			option.WithAPIKey(os.Getenv("LLM_API_KEY")),
			option.WithBaseURL(os.Getenv("LLM_API_URL")),
		),
	}
	if err := chatClient.initTools(ctx); err != nil {
		return nil, err
	}
	return chatClient, nil
}
func (c *Client) Chat(ctx context.Context, msg string) {
	fmt.Println(msg)
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(c.SystemPrompt),
			openai.SystemMessage(fmt.Sprintf("现在时间是：%s", time.Now().Format("2006-01-02 15:04:05"))),
			openai.UserMessage(msg),
		},
		Temperature: param.Opt[float64]{Value: 0.2},
		Model:       "qwen-plus",
		Tools:       c.ToolList,
	}

	for {
		var acc openai.ChatCompletionAccumulator
		var stream = c.LLMClient.Chat.Completions.NewStreaming(ctx, params)
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			if len(chunk.Choices) > 0 {
				// 即使在调用tool_call时也会有Content
				if chunk.Choices[0].Delta.Content != "" {
					fmt.Printf("\033[31m%s\033[0m", chunk.Choices[0].Delta.Content)
				}
				if openai.CompletionChoiceFinishReason(chunk.Choices[0].FinishReason) == openai.CompletionChoiceFinishReasonStop {
					return
				}
			}

			if _, ok := acc.JustFinishedContent(); ok {
				fmt.Println()
				fmt.Printf("finish-event: Content stream finished")
				fmt.Println()
			}

			if _, ok := acc.JustFinishedRefusal(); ok {
				fmt.Println()
				fmt.Printf("finish-event: refusal stream finished")
				fmt.Println()
			}

			if tool, ok := acc.JustFinishedToolCall(); ok {
				for _, choice := range acc.Choices {
					params.Messages = append(params.Messages, choice.Message.ToParam())
				}

				// TODO[Dokiy] 这里默认用户已经确认进行toolCall (2025/6/17)
				fmt.Println()
				fmt.Printf("call tool[%s]\n", tool.Name)
				fmt.Printf("tool.Arguments:%s\n", tool.Arguments)
				callToolResult, err := c.callTool(ctx, tool)
				if err != nil {
					return
				}

				for _, content := range callToolResult.Content {
					if textContext, ok := content.(mcp.TextContent); ok {
						fmt.Printf("tool.result: %s\n", strconv.Quote(textContext.Text))
						params.Messages = append(params.Messages, openai.ToolMessage(textContext.Text, tool.ID))
					}
				}
			}
		}
		if err := stream.Err(); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}
}
func (c *Client) initTools(ctx context.Context) error {
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.MCP.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		return fmt.Errorf("failed to initialize: %v", err)
	}

	var tools []mcp.Tool
	if serverInfo.Capabilities.Tools != nil {
		fmt.Println("Fetching available tools...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := c.MCP.ListTools(ctx, toolsRequest)
		if err != nil {
			fmt.Printf("Failed to list tools: %v", err)
		} else {
			tools = toolsResult.Tools
		}
	}

	c.ToolList = lo.Map(tools, func(tool mcp.Tool, index int) openai.ChatCompletionToolParam {
		return openai.ChatCompletionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name: tool.Name,
				Strict: param.Opt[bool]{
					Value: false,
				},
				Description: param.Opt[string]{
					Value: tool.Description,
				},
				Parameters: tool.InputSchema.Properties,
			},
			Type: "function",
		}
	})
	return nil
}
func (c *Client) callTool(ctx context.Context, tool openai.FinishedChatCompletionToolCall) (*mcp.CallToolResult, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(tool.Arguments), &args); err != nil {
		return nil, err
	}

	return c.MCP.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      tool.Name,
			Arguments: args,
		},
	})
}
