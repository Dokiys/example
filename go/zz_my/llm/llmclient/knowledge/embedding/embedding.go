package embedding

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

const (
	MaxChunkLength = 2048      // 最大分段字符长度
	ChunkOverlap   = 300       // 重叠字符长度
	ChunkSeparator = "---\n\n" // md分隔符
)

type Chunk struct {
	Text   string    `json:"text"`
	Vector []float64 `json:"vector"`
	Source string    `json:"source"`
}

type Embedding struct {
	Client openai.Client
	Model  string
}

type Callback interface {
	SearchTopK(ctx context.Context, queryVec []float64, topK int) ([]ScoredChunk, error)
}

func NewEmbedding(client openai.Client, model string) *Embedding {
	return &Embedding{Client: client, Model: model}
}

func (e *Embedding) Embedding(ctx context.Context, text string) ([]float64, error) {
	resp, err := e.Client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: param.Opt[string]{
				Value: text,
			},
		},
		Model:          e.Model,
		EncodingFormat: "float",
	})
	if err != nil {
		return nil, fmt.Errorf("调用接口向量化转换文本失败: %s", err)
	}

	// 因为这里用的OfString方式上传，理论上应该是中只返回一个元素
	if len(resp.Data) != 1 {
		panic(fmt.Sprintf("调用接口向量化转换文本返回Data数组长度不为1，resp: %v", resp))
	}
	return resp.Data[0].Embedding, nil
}

type ScoredChunk struct {
	Chunk `json:"chunk"`
	Score float64 `json:"score"`
}
