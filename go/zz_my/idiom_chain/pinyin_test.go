package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestPinyin(t *testing.T) {
	// 例如，查找匹配 abdy、a*b*y 这类 abbreviation（拼音首字母缩写）
	query := "fndp"
	matches, err := findIdiomsByPatternStreamed("idiom.json", query)
	if err != nil {
		panic(err)
	}

	for _, item := range matches {
		fmt.Printf("%s【%s】\n释义：%s\n\n", item.Word, item.Pinyin, item.Explanation)
	}
}

// findIdiomsByPatternStreamed 从大文件中匹配首字母模式并输出成语信息
func findIdiomsByPatternStreamed(filePath string, queryPattern string) ([]Idiom, error) {
	// 打开 JSON 文件
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer f.Close()

	// 准备正则表达式
	pattern := "^" + strings.ReplaceAll(queryPattern, "*", ".*") + "$"
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("正则编译失败: %w", err)
	}

	// 使用 JSON Decoder 流解析
	decoder := json.NewDecoder(f)

	// 预读取 `[` 开始标志
	tok, err := decoder.Token()
	if err != nil || tok != json.Delim('[') {
		return nil, fmt.Errorf("无效 JSON 格式，预期数组")
	}

	var result []Idiom

	// 循环每个成语对象
	for decoder.More() {
		var item Idiom
		if err := decoder.Decode(&item); err != nil {
			return nil, fmt.Errorf("成语解析失败: %w", err)
		}

		if re.MatchString(item.Abbr) {
			result = append(result, item)
		}
	}

	return result, nil
}
