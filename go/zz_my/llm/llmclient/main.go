package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"example/go/zz_my/llm/llmclient/knowledge/embedding"
	"example/go/zz_my/llm/llmclient/knowledge/memory"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

//go:embed system_prompt.md
var systemPrompt string

func main() {
	godotenv.Load()

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
	// c.Chat(ctx, "请帮我查询昨天核销对多的五个券批次，并列出券名、发放次数以及核销次数等信息。")
	c.Chat(ctx, "请在写一个成语接龙程序。运行程序后，有两种模式可以选择1.只要求接的新成语与上一个成语最后一个字同音，即拼音相同，或者是同一个字。2.要求接的成语必须和上一个成语是同一个字。用户选择后随机生成一个成语，然后等待用户输入新的成语，然后判断用户的输入是否符合要求，如果符合则按照要求随机找一个成语输出。如果不符合则输出用户接龙错误，要求重新输入。如此反复。直到无法接龙用户成语，时输出用户胜利，并推出。或者用户输入CTRL+C时退出。")
	// c.Chat(ctx, "为什么当请求biz__SqlExecutor工具获取到‘执行第1条SQL报错：sql: no rows in result set’报错的时候你会重复请求biz__SqlExecutor工具？")

	time.Sleep(time.Second * 2)
}
