package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	if err := c.RegisterMcpClient(ctx, []struct {
		Name string
		URL  string
	}{
		{
			Name: "local",
			URL:  "http://localhost:8080/mcp",
		},
		// {
		// 	Name: "biz",
		// 	URL:  "http://localhost:8081/coco/admin/mcp",
		// },
	}); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if err := c.InitTools(ctx); err != nil {
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
	McpList      map[string]*client.Client
	LLMClient    openai.Client
	ToolList     []openai.ChatCompletionToolParam
}

func NewClient(prompt string) (*Client, error) {
	chatClient := &Client{
		SystemPrompt: prompt,
		McpList:      make(map[string]*client.Client),
		LLMClient: openai.NewClient(
			option.WithAPIKey(os.Getenv("LLM_API_KEY")),
			option.WithBaseURL(os.Getenv("LLM_API_URL")),
		),
	}
	return chatClient, nil
}

func (c *Client) RegisterMcpClient(ctx context.Context, endPoints []struct {
	Name string
	URL  string
}) error {

	for _, point := range endPoints {
		httpTransport, err := transport.NewStreamableHTTP(point.URL)
		if err != nil {
			return fmt.Errorf("failed to create HTTP transport: %v", err)
		}

		mcpClient := client.NewClient(httpTransport)
		if err := mcpClient.Start(ctx); err != nil {
			return fmt.Errorf("%s", err)
		}
		c.McpList[point.Name] = mcpClient
	}

	return nil
}
func (c *Client) InitTools(ctx context.Context) error {
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	for name, mcpClient := range c.McpList {
		serverInfo, err := mcpClient.Initialize(ctx, mcp.InitializeRequest{})
		if err != nil {
			return fmt.Errorf("failed to initialize: %v", err)
		}

		var tools []mcp.Tool
		if serverInfo.Capabilities.Tools != nil {
			fmt.Println("Fetching available tools...")
			toolsRequest := mcp.ListToolsRequest{}
			toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
			if err != nil {
				fmt.Printf("Failed to list tools: %v", err)
			} else {
				tools = toolsResult.Tools
			}
		}

		c.ToolList = lo.Map(tools, func(tool mcp.Tool, index int) openai.ChatCompletionToolParam {
			fmt.Printf("Load tool[%s] from %s mcp client\n", tool.Name, name)
			return openai.ChatCompletionToolParam{
				Function: shared.FunctionDefinitionParam{
					Name: name + "__" + tool.Name,
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
		fmt.Println()
	}

	return nil
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

	var run = true
	for run {
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
					run = false
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
					fmt.Printf("Error: %s\n", err)
					return
				}

				for _, content := range callToolResult.Content {
					if textContext, ok := content.(mcp.TextContent); ok {
						fmt.Printf("tool.result: %s\n", strconv.Quote(textContext.Text))
						params.Messages = append(params.Messages, openai.ToolMessage(textContext.Text, tool.ID))
					}
				}
				fmt.Printf("finish-event: tool_call stream finished")
			}
		}
		if err := stream.Err(); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Println()
		fmt.Printf("\033[32mUsage: Total:%v CompletionTokens:%v, PromptTokens:%v\033[0m", acc.Usage.TotalTokens, acc.Usage.CompletionTokens, acc.Usage.PromptTokens)
	}
}
func (c *Client) callTool(ctx context.Context, tool openai.FinishedChatCompletionToolCall) (*mcp.CallToolResult, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(tool.Arguments), &args); err != nil {
		return nil, err
	}

	if s := strings.Split(tool.Name, "__"); len(s) == 2 {
		if mcpClient, ok := c.McpList[s[0]]; ok {
			return mcpClient.CallTool(ctx, mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      tool.Name,
					Arguments: args,
				},
			})
		}
	}

	return nil, fmt.Errorf("tool register error")
}
