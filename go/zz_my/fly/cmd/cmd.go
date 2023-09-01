package cmd

import (
	"reflect"

	"github.com/Dokiys/go_test/go/zz_my/fly/cmd/temp"
	"github.com/urfave/cli/v2"
)

type Commands struct {
	OneCmd  OneCmd
	TwoCmd  TwoCmd
	TempCmd temp.TempCmd
}

func NewCommands(cmd *Commands) cli.Commands {
	var cmds cli.Commands
	v := reflect.Indirect(reflect.ValueOf(cmd))
	ct := reflect.TypeOf(&cli.Command{})
	if v.Kind() != reflect.Struct {
		panic("reflect must be struct")
	}

	for i := 0; i < v.NumField(); i++ {
		cmd := v.Field(i).Convert(ct).Interface().(*cli.Command)
		cmds = append(cmds, cmd)
	}
	return cmds
}
