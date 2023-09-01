package temp

import (
	"fmt"

	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/urfave/cli/v2"
)

type SubTwoCmd *cli.Command

type SubTwo struct {
	Id string
}

func NewSubTwo(conf *conf.Config) *SubTwo {
	fmt.Println("init: SubTwo")
	return &SubTwo{conf.One}
}

func (o *SubTwo) Run() error {
	fmt.Println("sub:" + o.Id)
	return nil
}

func NewSubTwoCmd(conf *conf.Config) SubTwoCmd {
	return &cli.Command{
		Name:  "subtwo",
		Usage: "./app subtwo",
		Action: func(c *cli.Context) error {
			subtwo := NewSubTwo(conf)
			return subtwo.Run()
		},
	}
}
