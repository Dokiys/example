package temp

import (
	"fmt"

	"github.com/Dokiys/go_test/go/zz_my/fly"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/conf"
	"github.com/urfave/cli/v2"
)

type SubTwoCmd *cli.Command

type SubTwo struct {
	Id string
}

func NewSubTwo(conf *conf.Config) *SubTwo {
	fmt.Println("init: SubTwo")
	return &SubTwo{conf.Two}
}

func (o *SubTwo) Run() error {
	fmt.Println("sub:" + o.Id)
	return nil
}

func NewSubTwoCmd(fn fly.HandlerBuilder[*SubTwo]) *fly.Cmd {
	return &fly.Cmd{
		Name:   "subtwo",
		Usage:  "./app subtwo",
		Action: fn,
	}
}
