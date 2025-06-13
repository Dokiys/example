package mcp

import (
	"context"
	"fmt"
	"os"
	"testing"

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

func TestLLMClient(t *testing.T) {
	var ctx = context.Background()
	// 加载 .env 文件（默认当前目录下）
	godotenv.Load()

	r := openai.NewClient(
		option.WithAPIKey(os.Getenv("LLM_API_KEY")),
		option.WithBaseURL("https://aiproxy.verystar.net/v1/"),
	)

	c, err := NewMcpClient()
	if err != nil {
		t.Fatal(err)
	}

	tools, err := GetToolList(ctx, c)
	if err != nil {
		t.Fatal(err)
	}

	var stream = r.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
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
	})

	for stream.Next() {
		current := stream.Current()
		for _, choice := range current.Choices {
			if choice.FinishReason != "" {
				break
			}
			if choice.Delta.ToolCalls != nil {
				fmt.Printf("Call tools: %s\n", choice.Delta.ToolCalls[0].Function.Name)
				toolResult, err := c.CallTool(ctx, mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name:      choice.Delta.ToolCalls[0].Function.Name,
						Arguments: choice.Delta.ToolCalls[0].Function.Arguments,
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				fmt.Print(toolResult)
			}
			if choice.Delta.Content == "" {
				fmt.Print(choice.Delta.Content)
			}
		}
	}
}

func NewMcpClient() (*client.Client, error) {
	var ctx = context.Background()
	var httpURL = "http://localhost:8080/mcp"

	httpTransport, err := transport.NewStreamableHTTP(httpURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP transport: %v", err)
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
		return nil, fmt.Errorf("Failed to initialize: %v", err)
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
