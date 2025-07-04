package memory

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"sort"
	"strings"

	"example/go/zz_my/mcp/llmclient/knowledge/embedding"
)

//go:embed vector_data_gen.jsonl
var vectorDataGen embed.FS

func init() {
	// 初始化向量数据库
	if err := loadVectorDataset(); err != nil {
		panic(err)
	}
}

var VectorDataset = make([]embedding.Chunk, 0)

type EmbeddingCallback struct{}

func NewEmbeddingRecall() embedding.Callback {
	return &EmbeddingCallback{}
}

func loadVectorDataset() error {
	return fs.WalkDir(vectorDataGen, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".jsonl") {
			return nil
		}

		file, err := vectorDataGen.Open(path)
		if err != nil {
			return fmt.Errorf("打开文件失败: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 0

		for scanner.Scan() {
			line := scanner.Bytes()
			lineNumber++

			if len(line) == 0 {
				continue // 跳过空行
			}

			var chunk embedding.Chunk
			if err := json.Unmarshal(line, &chunk); err != nil {
				fmt.Printf("第%d行解析失败: %v\n", lineNumber, err)
				continue
			}

			VectorDataset = append(VectorDataset, chunk)
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("文件读取中出错: %w", err)
		}

		return nil
	})
}

// SearchTopK 在向量库中返回与 queryVec 最相似的 topK 个文本块
func (e *EmbeddingCallback) SearchTopK(_ context.Context, queryVec []float64, topK int) ([]embedding.ScoredChunk, error) {
	scored := make([]embedding.ScoredChunk, 0, len(VectorDataset))
	for _, c := range VectorDataset {
		score, err := cosineSimilarity(queryVec, c.Vector)
		if err == nil {
			scored = append(scored, embedding.ScoredChunk{Chunk: c, Score: score})
		}
	}

	// 按照分数排序（降序）
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// 取前 topK 的文本
	results := make([]embedding.ScoredChunk, 0, topK)
	for i := 0; i < topK && i < len(scored); i++ {
		results = append(results, scored[i])
	}

	return results, nil
}
func cosineSimilarity(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vector dimensions do not match")
	}

	dotProduct := 0.0
	normA := 0.0
	normB := 0.0
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("zero vector")
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}
