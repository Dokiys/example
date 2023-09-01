package cmd

import (
	"fmt"

	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/urfave/cli/v2"
)

type OneCmd *cli.Command

func NewOneCmd(conf *conf.Config) OneCmd {
	return &cli.Command{
		Name:  "one",
		Usage: "./app one",
		Action: func(c *cli.Context) error {
			fmt.Print(conf.One)
			return nil
		},
	}
}
