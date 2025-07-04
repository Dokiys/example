package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"example/go/zz_my/mcp/llmclient/knowledge/embedding"
	"example/go/zz_my/mcp/llmclient/knowledge/memory"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

//go:embed system_prompt.md
var systemPrompt string

func main() {
	var ctx = context.Background()
	c, err := NewClient(openai.NewClient(
		option.WithAPIKey(os.Getenv("LLM_API_KEY")),
		option.WithBaseURL(os.Getenv("LLM_API_URL")),
	),
		WithSystemPrompt(systemPrompt),
		WithEmbeddingClient(embedding.NewEmbedding(openai.NewClient(
			option.WithAPIKey(os.Getenv("LLM_API_KEY")),
			option.WithBaseURL(os.Getenv("LLM_EMBEDDING_API_URL")),
		), "qwen-text-embedding-v4")),
		WithEmbeddingRecall(memory.NewEmbeddingRecall()),
		// WithEmbeddingRecall(qdrant_search.NewEmbeddingRecall()),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if err := c.RegisterMcpClient(ctx, []EndPoint{
		// {
		// 	Name: "local",
		// 	URL:  "http://localhost:8080/mcp",
		// },
		{
			Name: "biz",
			URL:  "http://localhost:8081/coco/admin/mcp",
		},
	}); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if err := c.InitTools(ctx); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	c.Chat(ctx, "请帮我查询昨天核销对多的五个券批次，并列出券名、发放次数以及核销次数等信息。")
	// c.Chat(ctx, "为什么当请求biz__SqlExecutor工具获取到‘执行第1条SQL报错：sql: no rows in result set’报错的时候你会重复请求biz__SqlExecutor工具？")

	time.Sleep(time.Second * 2)
}
