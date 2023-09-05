package fly

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Handler interface {
	Run() error
}
type ConfigLoader interface {
	Load(v any)
}

type HandlerBuilder[T Handler] func(loader ConfigLoader) T

type actor interface {
	action(ConfigLoader) error
}

func (h HandlerBuilder[T]) action(loader ConfigLoader) error {
	cmd := h(loader)
	return cmd.Run()
}

type Cmd struct {
	Action actor

	// The name of the command
	Name string
	// A list of aliases for the command
	Aliases []string
	// A short description of the usage of this command
	Usage string
	// Custom text to show on USAGE section of help
	UsageText string
	// A longer explanation of how the command works
	Description string
	// A short description of the arguments of this command
	ArgsUsage string
	// The category the command is part of
	Category string
	// The function to call when checking for bash command completions
	BashComplete cli.BashCompleteFunc
	// An action to execute before any sub-subcommands are run, but after the context is ready
	// If a non-nil error is returned, no sub-subcommands are run
	Before cli.BeforeFunc
	// An action to execute after any subcommands are run, but after the subcommand has finished
	// It is run even if Action() panics
	After cli.AfterFunc
	// Execute this function if a usage error occurs.
	OnUsageError cli.OnUsageErrorFunc
	// List of child commands
	Subcommands []*Cmd
	// List of flags to parse
	Flags []cli.Flag
	// Treat all flags as normal arguments if true
	SkipFlagParsing bool
	// Boolean to hide built-in help command and help flag
	HideHelp bool
	// Boolean to hide built-in help command but keep help flag
	// Ignored if HideHelp is true.
	HideHelpCommand bool
	// Boolean to hide this command from help or completion
	Hidden bool
	// Boolean to enable short-option handling so user can combine several
	// single-character bool arguments into one
	// i.e. foobar -o -v -> foobar -ov
	UseShortOptionHandling bool
	// Full name of command for help, defaults to full command name, including parent commands.
	HelpName string
	// CustomHelpTemplate the text template for the command help topic.
	// cli.go uses text/template to render templates. You can
	// render custom help text by setting this variable.
	CustomHelpTemplate string
}

var commands []*Cmd

func Register(cmds ...*Cmd) {
	commands = append(commands, cmds...)
}

func Run(loader ConfigLoader) error {
	app := cli.NewApp()
	app.Name = "app"

	app.Commands = createCliCmds(loader, commands)
	return app.Run(os.Args)
}

func createCliCmds(loader ConfigLoader, cmds []*Cmd) (result []*cli.Command) {
	for _, cmd := range cmds {
		tmp := &cli.Command{
			Name:                   cmd.Name,
			Aliases:                cmd.Aliases,
			Usage:                  cmd.Usage,
			UsageText:              cmd.UsageText,
			Description:            cmd.Description,
			ArgsUsage:              cmd.ArgsUsage,
			Category:               cmd.Category,
			BashComplete:           cmd.BashComplete,
			Before:                 cmd.Before,
			After:                  cmd.After,
			Action:                 nil,
			OnUsageError:           cmd.OnUsageError,
			Subcommands:            createCliCmds(loader, cmd.Subcommands),
			Flags:                  cmd.Flags,
			SkipFlagParsing:        cmd.SkipFlagParsing,
			HideHelp:               cmd.HideHelp,
			HideHelpCommand:        cmd.HideHelpCommand,
			Hidden:                 cmd.Hidden,
			UseShortOptionHandling: cmd.UseShortOptionHandling,
			HelpName:               cmd.HelpName,
			CustomHelpTemplate:     cmd.CustomHelpTemplate,
		}

		if cmd.Action != nil {
			tmp.Action = func(ctx *cli.Context) error {
				if cmd.Action != nil {
					return cmd.Action.action(loader)
				}
				return nil
			}
		}
		result = append(result, tmp)
	}

	return result
}
