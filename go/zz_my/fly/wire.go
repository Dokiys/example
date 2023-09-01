//go:build wireinject

package main

import (
	"github.com/Dokiys/go_test/go/zz_my/fly/cmd"
	"github.com/Dokiys/go_test/go/zz_my/fly/cmd/temp"
	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

func initApp() cli.Commands {
	panic(wire.Build(wire.NewSet(
		conf.InitConf,

		temp.ProviderSet,

		cmd.NewOneCmd,

		cmd.NewTwoCmd,

		cmd.NewCommands,
		wire.Struct(new(cmd.Commands), "*"),
	)))
}
