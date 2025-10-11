package app

import (
	"os"

	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/log"
	"github.com/spf13/cobra"
)

type Command struct {
	use         string
	short       string
	discription string
	option      CliOption
	args        cobra.PositionalArgs
	cmd         *cobra.Command
	subCommands []*Command
	runfunc     CommandRunFunc
}

type CommandOption func(*Command)

func WithCommandShort(short string) CommandOption {
	return func(c *Command) {
		c.short = short
	}
}

func WithCommandDiscription(disc string) CommandOption {
	return func(c *Command) {
		c.discription = disc
	}
}

func WithCommandOption(option CliOption) CommandOption {
	return func(c *Command) {
		c.option = option
	}
}

type CommandRunFunc func() error

func WithCommandRunFunc(runfunc CommandRunFunc) CommandOption {
	return func(c *Command) {
		c.runfunc = runfunc
	}
}

func WithCommandPositionalArgs(args cobra.PositionalArgs) CommandOption {
	return func(c *Command) {
		c.args = args
	}
}

func NewCommand(use string, options ...CommandOption) *Command {
	cmd := &Command{
		use: use,
	}

	for _, o := range options {
		o(cmd)
	}

	return cmd
}

func (c *Command) buildCommand() {
	cmd := &cobra.Command{
		Use:   c.use,
		Short: c.short,
		Long:  c.discription,
		Args:  c.args,
	}
	cmd.SetErr(os.Stderr)
	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = true
	cmd.Flags().SetOutput(os.Stdout)
	cli.InitFS(cmd.Flags())

	if len(c.subCommands) > 0 {
		for _, sc := range c.subCommands {
			cmd.AddCommand(sc.CobraCommand())
		}
	}

	cmd.Run = c.runCommand

	nfs := c.option.Flags()

	for _, fs := range nfs.FlagSets {
		cmd.Flags().AddFlagSet(fs)
	}

	c.cmd = cmd
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cmd
}

func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if err := c.runfunc(); err != nil {
		log.Fatalf("cmd %s has error: %s\n", cmd.Use, err.Error())
	}
}

func (app *App) AddCommands(cmds ...*Command) {
	app.subCommands = append(app.subCommands, cmds...)
}

func (cmd *Command) AddCommands(cmds ...*Command) {
	cmd.subCommands = append(cmd.subCommands, cmds...)
}
