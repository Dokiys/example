package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"example/go/zz_my/clog"
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

	// TODO[Dokiy] to be continued! (2025/7/23)
	// 添加cursor的mcp.json导入支持
	if err := c.RegisterMcpClient(ctx, []EndPoint{
		// {
		// 	Name: "local",
		// 	URL:  "http://localhost:8080/mcp",
		// },
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

	var contextMsg [][]openai.ChatCompletionMessageParamUnion
	clog.Printf(clog.LevelCyan, "您好，我是您的智能助手，请问有什么可以帮您？\n")
	clog.Printf(clog.LevelGray, "请输入内容，输入 :end 结束，或 Ctrl+C 退出\n")

	// 捕捉 Ctrl+C 退出
	go exitSign()
	for {
		clog.Println()
		contextMsg = c.Chat(ctx, contextMsg, loadInput())
		clog.Println()
	}
}
func exitSign() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	fmt.Println("\nExit")
	os.Exit(0)
}

func loadInput() string {
	var input string
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		line = strings.TrimRight(line, "\r\n")
		if line == ":end" {
			break
		}
		input += line + "\n"
	}

	return input
}
