package cmd

import (
	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/Dokiys/go_test/go/zz_my/fly/server/one"
	"github.com/urfave/cli/v2"
)

type OneCmd *cli.Command

func NewOneCmd(conf *conf.Config) OneCmd {
	return &cli.Command{
		Name:  "one",
		Usage: "./app one",
		Action: func(c *cli.Context) error {
			oneSvc := one.NewOne(conf)
			return oneSvc.Run()
		},
	}
}
