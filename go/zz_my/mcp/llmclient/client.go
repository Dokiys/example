package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"example/go/zz_my/mcp/llmclient/embedding"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/samber/lo"
)

const embeddingPrompt = `
【参考资料】
%s

---

【问题】%s
`

type Client struct {
	SystemPrompt    string
	McpList         map[string]*client.Client
	LLMClient       openai.Client
	EmbeddingClient *embedding.Client
	TopK            int
	ToolList        []openai.ChatCompletionToolParam
}

type Option func(c *Client)

func NewClient(llmClient openai.Client, opts ...Option) (*Client, error) {
	chatClient := &Client{
		McpList:   make(map[string]*client.Client),
		LLMClient: llmClient,
		TopK:      3,
	}
	for _, opt := range opts {
		opt(chatClient)
	}

	return chatClient, nil
}
func WithSystemPrompt(prompt string) Option {
	return func(c *Client) {
		c.SystemPrompt = prompt
	}
}
func WithEmbeddingClient(client *embedding.Client) Option {
	return func(c *Client) {
		c.EmbeddingClient = client
	}
}

type EndPoint struct {
	Name string
	URL  string
}

func (c *Client) RegisterMcpClient(ctx context.Context, endPoints []EndPoint) error {
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
	if c.EmbeddingClient != nil {
		vec, err := c.EmbeddingClient.Embedding(ctx, msg)
		if err != nil {
			fmt.Printf("获取用户消息向量失败: %s\n", err)
		} else {
			chunks := embedding.SearchTopK(vec, embedding.VectorDataset, c.TopK)
			var reference = make([]string, 0, c.TopK)
			for i, chunk := range chunks {
				if chunk.Score <= 0 {
					// NOTE[Dokiy] (2025/7/2)
					// 可以设置相似性要求
				}
				reference = append(reference, fmt.Sprintf("%d. %s\n(来源：%s)", i+1, chunk.Text, chunk.Source))
			}

			msg = fmt.Sprintf(embeddingPrompt, strings.Join(reference, "\n\n"), msg)
			fmt.Printf("\033[37m知识库召回处理后：%s\033[0m\n", msg)
		}
	}
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(c.SystemPrompt),
			openai.SystemMessage(fmt.Sprintf("现在时间是：%s。请在调用工具的同时进行简单说明。", time.Now().Format("2006-01-02 15:04:05"))),
			openai.UserMessage(msg),
		},
		Temperature: param.Opt[float64]{Value: 0},
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
					fmt.Printf("\033[36m%s\033[0m", chunk.Choices[0].Delta.Content)
				}
				if openai.CompletionChoiceFinishReason(chunk.Choices[0].FinishReason) == openai.CompletionChoiceFinishReasonStop {
					fmt.Println()
					run = false
				}
			}

			if _, ok := acc.JustFinishedContent(); ok {
				// 如果LLM回复内容包含content和tool_call则会触发
				// fmt.Println()
				// fmt.Printf("finish-event: Content stream finished")
				// fmt.Println()
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

				fmt.Printf("ToolCall[%s]\n", tool.Name)
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
			}
		}
		if err := stream.Err(); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Printf("\033[37mUsage: Total:%v CompletionTokens:%v, PromptTokens:%v\033[0m\n", acc.Usage.TotalTokens, acc.Usage.CompletionTokens, acc.Usage.PromptTokens)
		fmt.Println()
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
					Name:      s[1],
					Arguments: args,
				},
			})
		}
	}

	return nil, fmt.Errorf("tool register error")
}
