package main

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func NewYY(vp *viper.Viper) *cli.Command {
	return &cli.Command{
		Name:  "yy",
		Usage: "意义",
		Action: func(c *cli.Context) error {
			contentsPrint(vp.GetStringSlice("fzsx.yy")...)

			return nil
		},
	}
}
func NewHXYY(vp *viper.Viper) *cli.Command {
	hxyy := vp.Get("fzsx.hxyy").([]interface{})
	cmds := []*cli.Command{}
	titles := make([]string, 0, len(hxyy))
	for i, e := range hxyy {
		m := e.(map[string]interface{})
		titles = append(titles, m["title"].(string))
		cmds = append(cmds, newHxyyCmd(i+1, m))
	}

	return &cli.Command{
		Name:        "hxyy",
		Usage:       "核心要义",
		Subcommands: cmds,
		Action: func(c *cli.Context) error {
			contentsPrint(titles...)

			return nil
		},
	}
}
func NewSJYQ(vp *viper.Viper) *cli.Command {
	sjyq := vp.Get("fzsx.sjyq").([]interface{})
	cmds := []*cli.Command{}
	titles := make([]string, 0, len(sjyq))
	for i, e := range sjyq {
		m := e.(map[string]interface{})
		titles = append(titles, m["title"].(string))
		cmds = append(cmds, newSjyqCmd(i+1, m))
	}

	return &cli.Command{
		Name:        "sjyq",
		Usage:       "实践要求",
		Subcommands: cmds,
		Action: func(c *cli.Context) error {
			contentsPrint(titles...)

			return nil
		},
	}
}

func newHxyyCmd(i int, m map[string]interface{}) *cli.Command {
	return &cli.Command{
		Name:  fmt.Sprint(i),
		Usage: fmt.Sprint(i),
		Subcommands: []*cli.Command{{
			Name: "content",
			Action: func(context *cli.Context) error {
				contentsPrint(toStringSlice(m["content"].([]interface{}))...)

				return nil
			},
		}},
		Action: func(context *cli.Context) error {
			contentsPrint(m["title"].(string) + ":")
			contentsPrint(m["what"].(string))
			contentsPrint(m["how"].(string))

			return nil
		},
	}
}
func newSjyqCmd(i int, m map[string]interface{}) *cli.Command {
	return &cli.Command{
		Name:  fmt.Sprint(i),
		Usage: fmt.Sprint(i),
		Action: func(context *cli.Context) error {
			contentsPrint(m["title"].(string) + ":")
			contentsPrint(toStringSlice(m["content"].([]interface{}))...)

			return nil
		},
	}
}

func toStringSlice(e []interface{}) []string {
	r := make([]string, 0, len(e))
	for _, v := range e {
		r = append(r, v.(string))
	}
	return r
}

func contentsPrint(contents ...string) {
	if len(contents) == 1 {
		fmt.Println(contents[0])
		return
	}

	for i, c := range contents {
		fmt.Printf("%d. %s\n", i+1, c)
	}
}

// TODO[Dokiy] 2023/2/20:
func splitLine(str string) {
	var buf bytes.Buffer
	var p int
	var isQuote = true
	var isNewline = false
	for i, r := range str {
		if (i+1)%30 == 0 {
			isNewline = !isNewline
		}
		if isNewline && !isQuote {
			buf.WriteString("\n")
		}

		if string(r) == "$" {
			p = i
			isQuote = !isQuote
		} else if string(r) == "}" {
			buf.WriteString(str[p:i])
			isQuote = !isQuote
			continue
		}
		buf.WriteRune(r)
	}

	fmt.Println(buf.String())
}
