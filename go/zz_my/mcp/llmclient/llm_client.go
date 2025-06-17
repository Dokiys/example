package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/samber/lo"
)

func main() {
	var ctx = context.Background()
	// 加载 .env 文件（默认当前目录下）
	_ = godotenv.Load("./.env")
	r := openai.NewClient(
		option.WithAPIKey(os.Getenv("LLM_API_KEY")),
		option.WithBaseURL(os.Getenv("LLM_API_URL")),
	)

	c, err := NewMcpClient()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	tools, err := GetToolList(ctx, c)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: param.Opt[string]{
							// Value: "你是谁？",
							Value: "请帮我查询openid为:yyds,且激活时间最晚的券码",
						},
						OfArrayOfContentParts: nil,
					},
				},
			},
		},
		Model: "qwen-plus",
		Tools: lo.Map(tools, func(tool mcp.Tool, index int) openai.ChatCompletionToolParam {
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
			}
		}),
	}

COMPLETION:
	for {
		var acc openai.ChatCompletionAccumulator
		var stream = r.Chat.Completions.NewStreaming(ctx, params)
		for stream.Next() {
			chunk := stream.Current()

			acc.AddChunk(chunk)
			if len(chunk.Choices) > 0 {
				if openai.CompletionChoiceFinishReason(chunk.Choices[0].FinishReason) == openai.CompletionChoiceFinishReasonStop {
					break COMPLETION
				}
				if chunk.Choices[0].Delta.Content != "" {
					print(chunk.Choices[0].Delta.Content)
				}
			}

			if content, ok := acc.JustFinishedContent(); ok {
				fmt.Printf("finish-event: Content stream finished: %s", content)
			}

			if refusal, ok := acc.JustFinishedRefusal(); ok {
				fmt.Printf("finish-event: refusal stream finished: %s", refusal)
			}

			if tool, ok := acc.JustFinishedToolCall(); ok {
				for _, choice := range acc.Choices {
					params.Messages = append(params.Messages, choice.Message.ToParam())
				}

				var args map[string]any
				if err := json.Unmarshal([]byte(tool.Arguments), &args); err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
				if desc, ok := args["description"]; ok {
					fmt.Printf("\n调用工具[%s]：%s\n", tool.Name, desc)
				}

				// TODO[Dokiy] 这里默认用户已经确认进行toolCall (2025/6/17)
				fmt.Printf("call tool: tool.Arguments:%s \n", tool.Arguments)
				callToolResult, err := c.CallTool(ctx, mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name:      tool.Name,
						Arguments: args,
					},
				})
				if err != nil {
					return
				}

				for _, content := range callToolResult.Content {
					if textContext, ok := content.(mcp.TextContent); ok {
						fmt.Printf("call tool: result: %s\n", strconv.Quote(textContext.Text))
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

func NewMcpClient() (*client.Client, error) {
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
	return c, nil
}

func GetToolList(ctx context.Context, c *client.Client) ([]mcp.Tool, error) {
	// Set up notification handler
	c.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification: %s\n", notification.Method)
	})

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize: %v", err)
	}

	var tools []mcp.Tool
	if serverInfo.Capabilities.Tools != nil {
		fmt.Println("Fetching available tools...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := c.ListTools(ctx, toolsRequest)
		if err != nil {
			fmt.Printf("Failed to list tools: %v", err)
		} else {
			tools = toolsResult.Tools
		}
	}

	return tools, nil
}
