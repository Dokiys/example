package main

import (
	"os"

	"github.com/Dokiys/go_test/go/zz_my/fly/conf"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "app"

	cmds, cleanUp := initApp(&conf.Config{One: "1", Two: "2"})
	defer cleanUp()
	app.Commands = cmds

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
