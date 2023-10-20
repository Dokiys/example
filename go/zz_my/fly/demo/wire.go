//go:build wireinject

package main

import (
	"github.com/Dokiys/go_test/go/zz_my/fly"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/server/one"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/server/two"
	"github.com/google/wire"
)

func NewOneFly(loader fly.ConfigLoader) *one.One {
	panic(wire.Build(providerSet))
}

func NewTwoFly(loader fly.ConfigLoader) *two.Two {
	panic(wire.Build(providerSet))
}
