package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func checkCommandExists(name string) {
	_, err := exec.LookPath(name)
	if err != nil {
		fmt.Printf("未找到 %s，请先安装。例如在 macOS 使用 brew install %s\n", name, name)
		os.Exit(2)
	}
}

func init() {
	// 检查依赖命令是否存在
	requiredCommands := []string{"git", "zsh", "fzf"}
	for _, cmd := range requiredCommands {
		checkCommandExists(cmd)
	}
}

// go build -o bin/aicommit
func main() {
	if os.Getenv("LLM_API_KEY") == "" || os.Getenv("LLM_API_URL") == "" {
		fmt.Println("Please set LLM_API_KEY and LLM_API_URL environment variable")
		os.Exit(2)
	}

	// git --no-pager diff -U0 b40e3e3eb ca2ec0240 ":(exclude)*.pb*.go"
	// # 1. 显示修改/删除等非新增文件的 diff
	// git --no-pager diff -U0--diff-filter=MDTR
	// # 2. 单独列出新增文件名
	// git --no-pager diff -U0 --name-only --diff-filter=A | sed 's/^/新增文件： /'
	cmdStr := `
git --no-pager diff -U0 --diff-filter=MDTR head ":(exclude)*.pb*.go";
echo;
git --no-pager diff -U0 --name-only --diff-filter=A head ":(exclude)*.pb*.go" | sed 's/^/新增文件： /'
`
	// 使用 zsh -c 执行一整个 shell 命令串
	cmd := exec.Command("zsh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("命令执行失败:", err)
		os.Exit(2)
	}

	var comment string
	var ctx = context.Background()
	var client = openai.NewClient(
		option.WithAPIKey(os.Getenv("LLM_API_KEY")),
		option.WithBaseURL(os.Getenv("LLM_API_URL")),
	)
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "qwen-plus-latest"
	}
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
请你作为一个摘要生成工具，帮我生成git commit提交时的comment信息。请根据以下文件差异，用简洁明了的完整句子（现在时）撰写一条有意义的 Git 提交信息，且长度不超过 74 个字符。
要求提交信息必须使用中文语言，开头必须加上以下之一的前缀：
| 类型      | 说明                                       |
|-----------|--------------------------------------------|
| feat    | 新增功能或特性                             |
| fix     | 修复 bug、问题                             |
| docs    | 修改文档内容（不影响代码逻辑）             |
| style   | 代码样式调整，不改变功能（空格、缩进等）  |
| refactor| 重构代码（不改变功能或修复 bug）           |
| perf    | 性能优化相关调整                           |
| test    | 添加或修改测试逻辑                         |
| chore   | 构建流程、脚手架、工具等变动               |
| build   | 构建相关改动（webpack、vite、脚本等）      |
| ci      | CI 配置或脚本变动（如 GitHub Actions）     |
| revert  | 撤销某次提交                               |
提交信息中包含的文件差异内容如下（删除的行以单个减号开头，新增的行以单个加号开头。如果没有提供任何变更内容则直接回复‘无内容变更’）：
`),
			openai.UserMessage(string(output)),
		},
		Model: model,
	})
	if err != nil {
		comment = fmt.Sprintln("LLM 处理失败:", err)
	} else {
		comment = resp.Choices[0].Message.Content
	}

	var selected bytes.Buffer
	fzf := exec.Command("fzf")
	commitOptions := []byte(comment + "\n") // 添加 AI 建议
	commitOptions = append(commitOptions, "fix: some\n"...)
	commitOptions = append(commitOptions, fmt.Sprintf("total tokens: %d\n", resp.Usage.TotalTokens)...)
	fzf.Stdin = bytes.NewReader(commitOptions)
	fzf.Stdout = &selected
	if err := fzf.Run(); err != nil {
		fmt.Println("取消选择或 fzf 出错:", err)
		return
	}

	fmt.Println(selected.String())
	os.Exit(0)
}
