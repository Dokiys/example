package temp

import (
	"github.com/Dokiys/go_test/go/zz_my/fly"
)

func NewTempCmd() *fly.Cmd {
	return &fly.Cmd{
		Name:  "temp",
		Usage: "./app temp",
		Subcommands: []*fly.Cmd{
			NewSubOneCmd(NewSubOneFly),
			NewSubTwoCmd(NewSubTwoFly),
		},
	}
}
