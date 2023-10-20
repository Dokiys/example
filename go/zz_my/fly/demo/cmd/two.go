package cmd

import (
	"github.com/Dokiys/go_test/go/zz_my/fly"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/server/two"
)

func NewTwoCmd(fn fly.HandlerBuilder[*two.Two]) *fly.Cmd {
	return &fly.Cmd{
		Name:   "two",
		Usage:  "./app two",
		Action: fn,
	}
}
