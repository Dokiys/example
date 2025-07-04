package qdrant_search

import (
	"context"
	"fmt"

	"example/go/zz_my/llm/llmclient/knowledge/embedding"
	"github.com/qdrant/go-client/qdrant"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

type EmbeddingCallback struct {
	CollectionName string
	QdClient       *qdrant.Client
}

func NewEmbeddingRecall() embedding.Callback {
	qdClient, err := qdrant.NewClient(&qdrant.Config{
		Host:                   "localhost",
		Port:                   6334,
		SkipCompatibilityCheck: true,
	})
	if err != nil {
		panic(fmt.Sprintf("创建qdrant链接失败: %s", err))
	}
	return &EmbeddingCallback{
		CollectionName: "qdrant",
		QdClient:       qdClient,
	}
}

func (e *EmbeddingCallback) SearchTopK(ctx context.Context, queryVec []float64, topK int) ([]embedding.ScoredChunk, error) {
	result, err := e.QdClient.Query(ctx, &qdrant.QueryPoints{
		CollectionName: e.CollectionName,
		Query:          qdrant.NewQueryDense(lo.Map(queryVec, func(v float64, _ int) float32 { return float32(v) })),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          proto.Uint64(uint64(topK)),
	})
	if err != nil {
		return nil, fmt.Errorf("dqrant查询数据失败：%s", err)
	}

	var topChunks []embedding.ScoredChunk
	for _, point := range result {
		topChunks = append(topChunks, embedding.ScoredChunk{
			Chunk: embedding.Chunk{
				Text:   point.Payload["text"].GetStringValue(),
				Vector: lo.Map(point.Vectors.GetVector().GetData(), func(v float32, _ int) float64 { return float64(v) }),
				Source: point.Payload["source"].GetStringValue(),
			},
			Score: float64(point.Score),
		})
	}

	return topChunks, nil
}
