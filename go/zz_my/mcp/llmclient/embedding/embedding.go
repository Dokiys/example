package embedding

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

//go:embed vector_data_gen/*
var vectorDataGen embed.FS

func init() {
	// 初始化向量数据库
	if err := LoadVectorDataset(vectorDataGen, "vector_data_gen"); err != nil {
		panic(err)
	}
}

const (
	MaxChunkLength = 2048      // 最大分段长度
	ChunkOverlap   = 300       // 重叠长度
	ChunkSeparator = "---\n\n" // md分隔符
)

type Chunk struct {
	Text   string    `json:"text"`
	Vector []float64 `json:"vector"`
	Source string    `json:"source"`
}

var VectorDataset = make([]Chunk, 0)

type Client struct {
	Client openai.Client
	Model  string
}

func NewClient(client openai.Client, model string) *Client {
	return &Client{Client: client, Model: model}
}

func (e *Client) Embedding(ctx context.Context, text string) ([]float64, error) {
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

func (e *Client) EmbeddingFile(ctx context.Context, rootFs embed.FS, root, outputDir string) error {
	return fs.WalkDir(rootFs, root, func(filepath string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		// 仅处理 md 文件
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		fmt.Printf("开始处理: %s\n", filepath)
		file, err := rootFs.Open(filepath)
		if err != nil {
			return fmt.Errorf("打开文件失败: %s, err: %w", filepath, err)
		}
		defer file.Close()

		base := path.Base(filepath)                       // xxx.md
		name := strings.TrimSuffix(base, path.Ext(base))  // xxx
		outputPath := path.Join(outputDir, name+".jsonl") // ${outputDir}/xxx.jsonl

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		enc := json.NewEncoder(outputFile)
		return splitTextWithOverlapReader(file, func(chunk string, index int) error {
			fmt.Printf("文件: %s，段 %d， 段落长度: %v:\n%s\n---\n", filepath, index, len(chunk), chunk)

			resp, err := e.Embedding(ctx, chunk)
			if err != nil {
				fmt.Println("处理失败:", err)
				return fmt.Errorf("调用接口向量化转换文本失败: %s", err)
			}

			if err := enc.Encode(Chunk{
				Text:   chunk,
				Vector: resp,
				Source: filepath,
			}); err != nil {
				return fmt.Errorf("向量化转换JSON序列化失败: %s", err)
			}
			return nil
		})

	})
}

func LoadVectorDataset(rootFs embed.FS, root string) error {
	return fs.WalkDir(rootFs, root, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".jsonl") {
			return nil
		}

		file, err := rootFs.Open(path)
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

			var chunk Chunk
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
func SearchTopK(queryVec []float64, chunks []Chunk, topK int) []ScoredChunk {
	scored := make([]ScoredChunk, 0, len(chunks))
	for _, c := range chunks {
		score, err := cosineSimilarity(queryVec, c.Vector)
		if err == nil {
			scored = append(scored, ScoredChunk{Chunk: c, Score: score})
		}
	}

	// 按照分数排序（降序）
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// 取前 topK 的文本
	results := make([]ScoredChunk, 0, topK)
	for i := 0; i < topK && i < len(scored); i++ {
		results = append(results, scored[i])
	}

	return results
}

var whitespaceRegex = regexp.MustCompile(`[\s\p{Zs}]+`)

// splitTextWithOverlapReader 基于 Reader 逐段读取文件
func splitTextWithOverlapReader(r io.Reader, onChunk func(chunk string, index int) error) error {
	var buf bytes.Buffer
	var prevTail string
	var index = 1
	var reader = bufio.NewReader(r)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				// 文件末尾处理最后一段
				finalText := prevTail + strings.TrimSpace(whitespaceRegex.ReplaceAllString(buf.String(), " "))
				if len(finalText) > 0 {
					if err := onChunk(finalText, index); err != nil {
						return fmt.Errorf("分段处理文件失败: %w", err)
					}
				}
				break
			}
			return err
		}

		// 当前段替换掉连续的空格、换行符和制表符后是否太长 or 遇到分割符？
		buf.WriteByte(b)
		trimChunk := strings.TrimSpace(whitespaceRegex.ReplaceAllString(buf.String(), " "))
		if len(prevTail)+len(trimChunk) >= MaxChunkLength || bytes.HasSuffix(buf.Bytes(), []byte(ChunkSeparator)) {
			fullText := prevTail + trimChunk

			if err := onChunk(fullText, index); err != nil {
				return fmt.Errorf("分段处理文件失败: %w", err)
			}
			index++

			// 保留末尾 ChunkOverlap 字符作为下一段的头部
			if len(fullText) > ChunkOverlap {
				prevTail = fullText[len(fullText)-ChunkOverlap:]
			} else {
				prevTail = fullText
			}

			buf.Reset()
		}
	}

	return nil
}

type ScoredChunk struct {
	Chunk `json:"chunk"`
	Score float64 `json:"score"`
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
