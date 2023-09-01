package temp

import (
	"fmt"

	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/urfave/cli/v2"
)

type SubOneCmd *cli.Command

type SubOne struct {
	Id string
}

func NewSubOne(conf *conf.Config) *SubOne {
	fmt.Println("init: SubOne")
	return &SubOne{conf.One}
}

func (o *SubOne) Run() error {
	fmt.Println("sub:" + o.Id)
	return nil
}

func NewSubOneCmd(conf *conf.Config) SubOneCmd {
	return &cli.Command{
		Name:  "subone",
		Usage: "./app subone",
		Action: func(c *cli.Context) error {
			subone := NewSubOne(conf)
			return subone.Run()
		},
	}
}
