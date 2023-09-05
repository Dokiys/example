package one

import (
	"fmt"

	"github.com/Dokiys/go_test/go/zz_my/fly/demo/conf"
)

type One struct {
	Id string
}

func NewOne(conf *conf.Config) *One {
	fmt.Println("init: One")
	return &One{conf.One}
}
func (o *One) Run() error {
	fmt.Println(o.Id)
	return nil
}
