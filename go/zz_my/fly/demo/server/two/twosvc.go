package two

import (
	"fmt"

	"github.com/Dokiys/go_test/go/zz_my/fly/demo/conf"
)

type Two struct {
	Id string
}

func NewTwo(conf *conf.Config) *Two {
	fmt.Println("init: Two")
	return &Two{conf.Two}
}
func (o *Two) Run() error {
	fmt.Println(o.Id)
	return nil
}
