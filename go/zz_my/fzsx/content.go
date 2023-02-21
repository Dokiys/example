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
		splitLine(contents[0])

		return
	}

	for i, c := range contents {
		splitLine(fmt.Sprintf("%d. %s\n", i+1, c))
	}
}

func splitLine(str string) {
	var buf bytes.Buffer
	var p int
	var isQuote = false
	var isNewline = false
	for i, r := range str {
		if (i/3+1)%30 == 0 && !isNewline {
			isNewline = true
		}

		if string(r) == "$" {
			p = i
			isQuote = true
		} else if string(r) == "}" {
			buf.WriteString(str[p : i+1])
			isQuote = false
		} else if !isQuote {
			buf.WriteRune(r)
		}

		if !isQuote && isNewline {
			buf.WriteString("\n")
			isNewline = false
		}

	}

	fmt.Println(buf.String())
}
