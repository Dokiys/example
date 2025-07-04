package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mozillazg/go-pinyin"
)

type Idiom struct {
	Word        string `json:"word"`
	Abbr        string `json:"abbreviation"`
	Pinyin      string `json:"pinyin"`
	Explanation string `json:"explanation"`
}

// LoadAllIdioms 加载全部成语（可选缓存）
func LoadAllIdioms(filePath string) ([]Idiom, error) {
	f, err := idiomsJson.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	tok, err := decoder.Token()
	if err != nil || tok != json.Delim('[') {
		return nil, fmt.Errorf("不是 JSON 数组")
	}

	var idioms []Idiom
	for decoder.More() {
		var idiom Idiom
		if err := decoder.Decode(&idiom); err != nil {
			return nil, err
		}
		// 忽略不合法词
		if len([]rune(idiom.Word)) != 4 {
			continue
		}
		idioms = append(idioms, idiom)
	}

	return idioms, nil
}

// 获取汉字首字或尾字的拼音（不含声调）
func getCharPinyin(r rune) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	p := pinyin.SinglePinyin(r, args)
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

// 随机选择一条接龙成语
func findNextIdiom(word string, idioms []Idiom, mode int) *Idiom {
	lastRune := []rune(word)[len([]rune(word))-1]
	lastPinyin := getCharPinyin(lastRune)

	rand.Shuffle(len(idioms), func(i, j int) {
		idioms[i], idioms[j] = idioms[j], idioms[i]
	})

	for _, idiom := range idioms {
		firstRune := []rune(idiom.Word)[0]
		if mode == 2 {
			if firstRune == lastRune {
				return &idiom
			}
		} else {
			firstPinyin := getCharPinyin(firstRune)
			if firstRune == lastRune || firstPinyin == lastPinyin {
				return &idiom
			}
		}
	}
	return nil
}

// 用户输入的成语是否合法，并且能接上上一个词
func isValidLink(user string, last string, idiomsMap map[string]Idiom, mode int) bool {
	id, ok := idiomsMap[user]
	if !ok || len([]rune(id.Word)) != 4 {
		fmt.Println("不是有效的成语")
		return false
	}

	firstRune := []rune(user)[0]
	lastRune := []rune(last)[len([]rune(last))-1]

	if mode == 2 {
		return firstRune == lastRune
	}

	return firstRune == lastRune || getCharPinyin(firstRune) == getCharPinyin(lastRune)
}

//go:embed idiom.json
var idiomsJson embed.FS

func main() {
	idioms, err := LoadAllIdioms("idiom.json")
	if err != nil {
		panic(err)
	}
	idiomMap := make(map[string]Idiom)
	for _, id := range idioms {
		idiomMap[id.Word] = id
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("✅ 成语接龙开始！请选择接龙模式：")
	fmt.Println("1. 拼音或字相同")
	fmt.Println("2. 字必须相同")
	fmt.Println("输入 1 或 2：")

	var mode int
	for {
		fmt.Fscanln(reader, &mode)
		if mode == 1 || mode == 2 {
			break
		}
		fmt.Println("请输入有效模式 (1或2)：")
	}

	// 捕捉 Ctrl+C 退出
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		fmt.Println("\n👋 已手动退出游戏")
		os.Exit(0)
	}()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	current := idioms[r.Intn(len(idioms))]

	fmt.Printf("🤖 开始接龙：%s (%s)\n", current.Word, current.Pinyin)
	fmt.Println("你来接（输入成语）：")

	var hinted bool
	var textCh = make(chan string)
	go func() {
		for {
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			textCh <- input
		}
	}()
	for {
		var text string
		select {
		case <-time.After(time.Second * 10):
			if !hinted {
				fmt.Println("\033[37m🤖 被难住了可以输入'提示'哦！\033[0m")
				hinted = true
			}
			continue
		case text = <-textCh:
		}

		if text == "提示" || strings.ToLower(text) == "hint" {
			tip := findNextIdiom(current.Word, idioms, mode)
			if tip == nil {
				fmt.Println("🎉 太棒了！我也想不到，你赢了！")
				break
			}
			fmt.Printf("🤖 我帮你接：%s (%s)\n释义：%s\n\n", tip.Word, tip.Pinyin, tip.Explanation)
			fmt.Println("你来接（输入成语）：")
			current = *tip
			continue
		}

		if !isValidLink(text, current.Word, idiomMap, mode) {
			fmt.Println("❌ 不符合接龙规则，请重新输入！")
			continue
		}

		next := findNextIdiom(text, idioms, mode)
		if next == nil {
			fmt.Println("🎉 太棒了！我接不下去了，你赢了！")
			break
		}

		fmt.Printf("🤖 我接：%s (%s)\n解释：%s\n\n", next.Word, next.Pinyin, next.Explanation)
		fmt.Println("你来接（输入成语）：")
		current = *next
	}
}
