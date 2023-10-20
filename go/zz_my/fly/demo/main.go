package main

import (
	"github.com/Dokiys/go_test/go/zz_my/fly"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/cmd"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/cmd/temp"
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/conf"
)

func main() {
	fly.Register(
		cmd.NewOneCmd(NewOneFly),
		cmd.NewTwoCmd(NewTwoFly),
		temp.NewTempCmd(),
	)

	err := fly.Run(&conf.BootConfigLoader{})
	if err != nil {
		panic(err)
	}
}
