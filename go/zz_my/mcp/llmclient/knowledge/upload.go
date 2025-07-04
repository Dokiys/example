package knowledge

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"strings"

	"example/go/zz_my/mcp/llmclient/knowledge/embedding"
)

//go:embed asset/*
var asset embed.FS

func UploadKnowledge(ctx context.Context, e *embedding.Embedding, fn func(chunk embedding.Chunk) error) error {
	return fs.WalkDir(asset, "asset", func(filepath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 仅处理 md 文件
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		fmt.Printf("开始处理: %s\n", filepath)
		file, err := asset.Open(filepath)
		if err != nil {
			return fmt.Errorf("打开文件失败: %s, err: %w", filepath, err)
		}
		defer file.Close()

		return splitTextWithOverlapReader(file, func(text string, index int) error {
			fmt.Printf("文件: %s，段 %d， 段落长度: %v:\n%s\n---\n", filepath, index, len(text), text)

			resp, err := e.Embedding(ctx, text)
			if err != nil {
				fmt.Println("处理失败:", err)
				return fmt.Errorf("调用接口向量化转换文本失败: %s", err)
			}

			chunk := embedding.Chunk{
				Text:   text,
				Vector: resp,
				Source: filepath,
			}
			if err := fn(chunk); err != nil {
				return err
			}

			return nil
		})
	})
}

var whitespaceRegex = regexp.MustCompile(`[\s\p{Zs}]+`)

// splitTextWithOverlapReader 基于 Reader 逐段读取文件
func splitTextWithOverlapReader(r io.Reader, onChunk func(chunk string, index int) error) error {
	var buf strings.Builder
	var prevTail string
	var index = 1
	var reader = bufio.NewReader(r)

	for {
		// NOTE[Dokiy] (2025/7/4)
		// 文本文件按照byte读取会把utf-8字符截断，所以需要用rune读取
		rn, _, err := reader.ReadRune()
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
		buf.WriteRune(rn)
		trimChunk := strings.TrimSpace(whitespaceRegex.ReplaceAllString(buf.String(), " "))
		if len(prevTail)+len(trimChunk) >= embedding.MaxChunkLength || bytes.HasSuffix([]byte(buf.String()), []byte(embedding.ChunkSeparator)) {
			fullText := prevTail + trimChunk

			if err := onChunk(fullText, index); err != nil {
				return fmt.Errorf("分段处理文件失败: %w", err)
			}
			index++

			// 截取末尾部分以便作为下一段的开头（UTF-8 安全，按字符取）
			runes := []rune(fullText)
			if len(runes) > embedding.ChunkOverlap {
				prevTail = string(runes[len(runes)-embedding.ChunkOverlap:])
			} else {
				prevTail = fullText
			}

			buf.Reset()
		}
	}

	return nil
}
