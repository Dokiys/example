package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "app"

	cmds := initApp()
	app.Commands = cmds

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
