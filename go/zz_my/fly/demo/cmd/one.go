package cmd

import (
	"github.com/Dokiys/go_test/go/zz_my/fly"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/server/one"
)

func NewOneCmd(fn fly.HandlerBuilder[*one.One]) *fly.Cmd {
	return &fly.Cmd{
		Name:   "one",
		Usage:  "./app one",
		Action: fn,
	}
}
