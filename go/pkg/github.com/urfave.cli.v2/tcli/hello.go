package tcli

import (
	"fmt"
)

func NewHelloCmd() *cli.Command {
	return &cli.Command{
		Name:        "hello",
		Usage:       "Say Hello",
		UsageText:   "123",
		Description: "say hello",
		ArgsUsage:   "ArgsUsage",
		Flags:       []cli.Flag{},
		Action: func(context *cli.Context) error {
			fmt.Println("Hello!")
			return nil
		},
	}
}
