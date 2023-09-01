package temp

import (
	"reflect"

	"github.com/urfave/cli/v2"
)

type Temp struct {
	SubOneCmd SubOneCmd
	SubTwoCmd SubTwoCmd
}

type TempCmd *cli.Command

func NewTempCmd(temp *Temp) TempCmd {
	return &cli.Command{
		Name:        "temp",
		Usage:       "./app temp",
		Subcommands: newSubCommands(temp),
	}
}

func newSubCommands(cmd *Temp) cli.Commands {
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
