package embedding

import (
	"context"
	"embed"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

//go:embed asset/*
var asset embed.FS

func TestEmbeddingAsset(t *testing.T) {
	godotenv.Load("../.env")

	ctx := context.Background()
	client := NewClient(openai.NewClient(
		option.WithAPIKey(os.Getenv("LLM_API_KEY")),
		option.WithBaseURL(os.Getenv("LLM_EMBEDDING_API_URL")),
	), "qwen-text-embedding-v4")
	if err := client.EmbeddingFile(ctx, asset, "asset", "vector_data_gen"); err != nil {
		t.Fatal(err)
	}
}

func TestLoadVectorDataset(t *testing.T) {
	godotenv.Load("../.env")

	ctx := context.Background()
	client := Client{
		Model: "qwen-text-embedding-v4",
		Client: openai.NewClient(
			option.WithAPIKey(os.Getenv("LLM_API_KEY")),
			option.WithBaseURL(os.Getenv("LLM_EMBEDDING_API_URL")),
		),
	}

	if err := LoadVectorDataset(vectorDataGen, "vector_data_gen"); err != nil {
		t.Fatal(err)
	}

	// 从通义千问 embedding 得到 query 向量
	queryVec, err := client.Embedding(ctx, "帮我查一下最近发出的两张电子券码的详情信息")
	if err != nil {
		t.Fatal(err)
	}

	topChunks := SearchTopK(queryVec, VectorDataset, 2)
	t.Log(topChunks)
}
