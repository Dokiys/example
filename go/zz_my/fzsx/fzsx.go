package main

import (
	"bytes"
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

//go:embed fzsx.yml
var FS embed.FS

func main() {
	vp := viper.New()
	vp.SetConfigType("yaml")
	bs, err := FS.ReadFile("fzsx.yml")
	if err != nil {
		panic(err)
	}

	if err := vp.ReadConfig(bytes.NewBuffer(bs)); err != nil {
		panic(err)
	}

	Run(func(signal os.Signal) {}, []*cli.Command{
		NewYY(vp),
		NewHXYY(vp),
		NewSJYQ(vp),
	}...)
}

const VERSION = "v1.0.0"

var notifySignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

func waitSign(stop chan<- struct{}, cleanup func(os.Signal)) {
	n := make(chan os.Signal)
	signal.Notify(n, notifySignals...)

	cleanup(<-n)
	stop <- struct{}{}
}

func Run(cleanup func(os.Signal), cmds ...*cli.Command) {
	var stop = make(chan struct{})

	app := cli.NewApp()
	app.Name = "fzsx"
	app.Usage = "fzsx [COMMANDS]"
	app.UsageText = "fzsx [global options] command [command options] [arguments...]"
	app.Copyright = "Copyright ©2023 Dokiy"
	app.Version = VERSION
	app.Authors = []*cli.Author{{
		Name:  "Dokiy",
		Email: "",
	}}
	app.Commands = cmds
	app.After = func(c *cli.Context) error {
		close(stop)
		return nil
	}

	go waitSign(stop, cleanup)
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
	<-stop
}
