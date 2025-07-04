package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"example/go/zz_my/llm/llmclient/knowledge/embedding"
	"example/go/zz_my/llm/llmclient/knowledge/memory"
	"example/go/zz_my/llm/llmclient/knowledge/qdrant_search"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/qdrant/go-client/qdrant"
	"github.com/samber/lo"
)

func TestEmbeddingMemory(t *testing.T) {
	godotenv.Load("../.env")

	ctx := context.Background()
	client := embedding.NewEmbedding(openai.NewClient(
		option.WithAPIKey(os.Getenv("LLM_API_KEY")),
		option.WithBaseURL(os.Getenv("LLM_EMBEDDING_API_URL")),
	), "qwen-text-embedding-v4")

	outputPath := path.Join("./memory/vector_data_gen.jsonl")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(outputFile)
	if err := UploadKnowledge(ctx, client, func(chunk embedding.Chunk) error {
		if err := enc.Encode(chunk); err != nil {
			return fmt.Errorf("向量化转换JSON序列化失败: %s", err)
		}
		return nil

	}); err != nil {
		t.Fatal(err)
	}
}
func TestEmbeddingQdrant(t *testing.T) {
	godotenv.Load("../.env")

	ctx := context.Background()
	client := embedding.NewEmbedding(openai.NewClient(
		option.WithAPIKey(os.Getenv("LLM_API_KEY")),
		option.WithBaseURL(os.Getenv("LLM_EMBEDDING_API_URL")),
	), "qwen-text-embedding-v4")

	qdClient, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		return
	}

	var collectionName = "qdrant"
	if err := UploadKnowledge(ctx, client, func(chunk embedding.Chunk) error {
		_, err := qdClient.Upsert(context.Background(), &qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points: []*qdrant.PointStruct{
				{
					Id:      qdrant.NewIDNum(djbHash64([]byte(chunk.Text))),
					Vectors: qdrant.NewVectorsDense(lo.Map(chunk.Vector, func(v float64, _ int) float32 { return float32(v) })),
					Payload: qdrant.NewValueMap(map[string]any{"text": chunk.Text, "source": chunk.Source}),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("插入qdrant向量数据失败：%s", err)
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
func TestMemorySearchTopK(t *testing.T) {
	godotenv.Load("../.env")

	ctx := context.Background()
	client := embedding.Embedding{
		Model: "qwen-text-embedding-v4",
		Client: openai.NewClient(
			option.WithAPIKey(os.Getenv("LLM_API_KEY")),
			option.WithBaseURL(os.Getenv("LLM_EMBEDDING_API_URL")),
		),
	}

	// 从通义千问 embedding 得到 query 向量
	queryVec, err := client.Embedding(ctx, "帮我查一下最近发出的两张电子券码的详情信息")
	if err != nil {
		t.Fatal(err)
	}

	topChunks, err := memory.NewEmbeddingRecall().SearchTopK(ctx, queryVec, 3)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(topChunks)
}

func TestMemoryQdrantSearch(t *testing.T) {
	godotenv.Load("../.env")

	ctx := context.Background()
	client := embedding.Embedding{
		Model: "qwen-text-embedding-v4",
		Client: openai.NewClient(
			option.WithAPIKey(os.Getenv("LLM_API_KEY")),
			option.WithBaseURL(os.Getenv("LLM_EMBEDDING_API_URL")),
		),
	}

	var topChunks []embedding.ScoredChunk
	// 从通义千问 embedding 得到 query 向量
	queryVec, err := client.Embedding(ctx, "帮我查一下最近发出的两张电子券码的详情信息")
	if err != nil {
		t.Fatal(err)
	}
	topChunks, err = qdrant_search.NewEmbeddingRecall().SearchTopK(ctx, queryVec, 3)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(topChunks)
}

func djbHash64(str []byte) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(str); i++ {
		hash = (hash<<5 + hash) + uint64(str[i])
	}
	return hash
}
