package temp

import (
	"fmt"

	"github.com/Dokiys/go_test/go/zz_my/fly"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/conf"
)

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

func NewSubOneCmd(fn fly.HandlerBuilder[*SubOne]) *fly.Cmd {
	return &fly.Cmd{
		Name:   "subone",
		Usage:  "./app subone",
		Action: fn,
	}
}
