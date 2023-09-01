package cmd

import (
	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/Dokiys/go_test/go/zz_my/fly/server/two"
	"github.com/urfave/cli/v2"
)

type TwoCmd *cli.Command

func NewTwoCmd(conf *conf.Config) TwoCmd {
	return &cli.Command{
		Name:  "two",
		Usage: "./app two",
		Action: func(c *cli.Context) error {
			twoSvc := two.NewTwo(conf)
			return twoSvc.Run()
		},
	}
}
