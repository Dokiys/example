package src

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestPlaceholder(t *testing.T) {
	zhangsan := struct {
		Name string
	}{Name: "zhangsan"}

	t.Logf("%-s | %v", fmtHan(30, "相应值的默认格式"), zhangsan)
	t.Logf("%-s | %+v", fmtHan(30, "打印结构体时添加字段名"), zhangsan)
	t.Logf("%-s | %#v", fmtHan(30, "相应值的Go语法表示"), zhangsan)
	t.Logf("%-s | %T", fmtHan(30, "相应值的类型的Go语法表示"), zhangsan)
	t.Logf("%-s | %p", fmtHan(30, "十六进制表示，前缀 0x"), &zhangsan)
	t.Logf("%-s | %%", fmtHan(30, "转义%"))
	t.Log()
	t.Logf("%-s | %b", fmtHan(30, "二进制表示"), 10)
	t.Logf("%-s | %c", fmtHan(30, "相应Unicode码所表示的字符"), 0x4E2D)
	t.Logf("%-s | %d", fmtHan(30, "十进制表示"), 0x12)
	t.Logf("%-s | %05d", fmtHan(30, "十进制表示（0填充到至少5位）"), 10)
	t.Logf("%-s | %o", fmtHan(30, "八进制表示"), 10)
	t.Logf("%-s | %q", fmtHan(30, "单引号围绕的字符"), 0x4E2D)
	t.Logf("%-s | %x", fmtHan(30, "十六进制表示，字母为小写"), 13)
	t.Logf("%-s | %X", fmtHan(30, "十六进制表示，字母为大写"), 13)
	t.Logf("%-s | %U", fmtHan(30, "Unicode格式：U+1234"), 0x4E2D)
	t.Log()
	t.Logf("%-s | %e", fmtHan(30, "科学计数法"), 13.14)
	t.Logf("%-s | %E", fmtHan(30, "科学计数法"), 13.14)
	t.Logf("%-s | %f", fmtHan(30, "包括小数"), 13.14)
	t.Logf("%-s | %g", fmtHan(30, "自动选择最短的表示形式"), math.Pi)
	t.Logf("%-s | %G", fmtHan(30, "自动选择最短的表示形式"), math.Pi)
	t.Log()
	t.Logf("%-s | %s", fmtHan(30, "输出字符串"), "中文")
	t.Logf("%-s | %s", fmtHan(30, "输出[]byte"), []byte("中文"))
	t.Logf("%-s | %15s", fmtHan(30, "填充字符串至少到15位"), "中文")
	t.Logf("%-s | %q", fmtHan(30, "双引号围绕的字符串"), "中文")
	t.Logf("%-s | %x", fmtHan(30, "十六进制小写字母表示"), "中文")
	t.Logf("%-s | %X", fmtHan(30, "十六进制大写字母表示"), "中文")
}

func TestColor(t *testing.T) {
	color_table := map[string]string{
		"Reset": "\x1b[0m",

		"Bright":     "\x1b[1m",
		"Dim":        "\x1b[2m",
		"Underscore": "\x1b[4m",
		"Blink":      "\x1b[5m",
		"Reverse":    "\x1b[7m",
		"Hidden":     "\x1b[8m",

		"FgBlack":   "\x1b[30m",
		"FgRed":     "\x1b[31m",
		"FgGreen":   "\x1b[32m",
		"FgYellow":  "\x1b[33m",
		"FgBlue":    "\x1b[34m",
		"FgMagenta": "\x1b[35m",
		"FgCyan":    "\x1b[36m",
		"FgWhite":   "\x1b[37m",

		"BgBlack":   "\x1b[40m",
		"BgRed":     "\x1b[41m",
		"BgGreen":   "\x1b[42m",
		"BgYellow":  "\x1b[43m",
		"BgBlue":    "\x1b[44m",
		"BgMagenta": "\x1b[45m",
		"BgCyan":    "\x1b[46m",
		"BgWhite":   "\x1b[47m",
	}

	// fg
	for k, v := range color_table {
		if strings.HasPrefix(k, "Fg") {
			t.Logf("%15s %s", k, fmt.Sprintf("%s%s\u001B[0m", v, k))
		}
	}

	t.Log()

	// bg
	for k, v := range color_table {
		if strings.HasPrefix(k, "Bg") {
			t.Logf("%15s %s", k, fmt.Sprintf("%s%s\u001B[0m", v, k))
		}
	}

	t.Log()

	// misc
	for k, v := range color_table {
		if !(strings.HasPrefix(k, "Fg") || strings.HasPrefix(k, "Bg")) {
			t.Logf("%15s %s", k, fmt.Sprintf("%s%s\u001B[0m", v, k))
		}
	}
}
