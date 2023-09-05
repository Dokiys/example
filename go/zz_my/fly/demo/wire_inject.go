//go:build wireinject

package main

import (
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/cmd/temp"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/conf"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/server/one"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/server/two"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	conf.NewConfig,
	one.NewOne,
	two.NewTwo,
	temp.NewTempCmd,
)
