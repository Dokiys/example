//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Dokiys/go_test/go/zz_my/fly/cmd"
	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

var providerSet = wire.NewSet(
	cmd.NewOneCmd,
	cmd.NewTwoCmd,
	cmd.NewCommands,
	wire.Struct(new(cmd.Commands), "*"),
)

func initApp(conf *conf.Config) (cli.Commands, func()) {
	panic(wire.Build(providerSet))
}
