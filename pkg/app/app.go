package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/version"
	"github.com/neee333ko/component-base/pkg/version/verflag"

	"github.com/neee333ko/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var progressMessage string = color.GreenString("==>")

func init() {
	log.Init(log.InitOptions(log.WithFormat("console")))
}

// App is used to build the root cmd for application
type App struct {
	name        string
	basename    string
	short       string
	discription string
	option      CliOption
	noVersion   bool
	noConfig    bool
	silence     bool
	runFunc     RunFunc
	args        cobra.PositionalArgs
	cmd         *cobra.Command
	subCommands []*Command
}

type Option func(*App)

func WithShort(short string) Option {
	return func(a *App) {
		a.short = short
	}
}

func WithDiscription(disc string) Option {
	return func(a *App) {
		a.discription = disc
	}
}

func WithOption(option CliOption) Option {
	return func(a *App) {
		a.option = option
	}
}

func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

type RunFunc func(string)

func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

func WithSubCommand(subCmds []*Command) Option {
	return func(a *App) {
		a.subCommands = subCmds
	}
}

func WithPositionalArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

func WithDefaultPositionalArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("must not have args")
				}
			}

			return nil
		}
	}
}

// New returns a new app with name, basename, and Option.
func New(name string, basename string, options ...Option) *App {
	app := &App{
		name:     name,
		basename: FormatName(basename),
	}

	for _, o := range options {
		o(app)
	}

	// app.buildCommand()

	return app
}

func (app *App) buildCommand() {
	cmd := &cobra.Command{
		Use:           app.basename,
		Short:         app.short,
		Long:          app.discription,
		Args:          app.args,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	cmd.Flags().SetOutput(os.Stdout)
	cli.InitFS(cmd.Flags())

	if len(app.subCommands) > 0 {
		for _, sc := range app.subCommands {
			cmd.AddCommand(sc.CobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(app.basename))
	}

	cmd.Run = app.runCommand

	var nfs cli.NamedFlagSets

	for _, fs := range app.option.Flags().FlagSets {
		nfs.AddFlagSet("global", fs)
	}

	if !app.noConfig {
		verflag.AddFlag(nfs.FlagSets["global"], app.basename)
	}

	if !app.noVersion {
		AddConfigFlag(nfs.FlagSets["global"], app.basename)
	}

	cli.AddHelpFlag(nfs.FlagSets["global"], app.basename)

	cmd.Flags().AddFlagSet(nfs.FlagSets["global"])

	app.cmd = cmd
}

func (app *App) buildApp() {
	for _, sc := range app.subCommands {
		build(sc)
	}

	app.buildCommand()
}

func build(cmd *Command) {
	for _, sc := range cmd.subCommands {
		build(sc)
	}

	cmd.buildCommand()
}

func (app *App) Run() {
	app.buildApp()

	if err := app.cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func (app *App) runCommand(cmd *cobra.Command, args []string) {
	printcwd()

	if !app.noVersion {
		verflag.PrintAndExit(app.cmd.Flags())
	}

	if !app.noConfig {
		if err := viper.BindPFlags(app.cmd.Flags()); err != nil {
			log.Fatal(err.Error())
		}

		if err := viper.Unmarshal(app.option); err != nil {
			log.Fatal(err.Error())
		}
	}

	if !app.silence {
		if !app.noVersion {
			log.Info(fmt.Sprintf("\n%s", version.Get().String()))
		}

		if !app.noConfig {
			log.Info(printConfigFlag())
		}
	}

	log.Infof("\n%s", cli.PrintFlagSet(cmd.Flags()))

	app.applyOption()

	app.runFunc(app.basename)
}

func (app *App) applyOption() {
	if co, ok := app.option.(CompleteOption); ok {
		err := co.Complete()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if err := app.option.Validate(); err != nil {
		log.Fatal(err.ToAggregate().Error())
	}

	if po, ok := app.option.(PrintOption); ok && !app.silence {
		log.Info(po.String())
	}
}

// FormatName returns binary executable file name when OS is WINDOWS
func FormatName(basename string) string {
	if runtime.GOOS == "windows" {
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return strings.ToLower(basename)
}

func printcwd() {
	path, err := os.Getwd()
	if err != nil {
		log.Info(fmt.Sprintf("Working Directory: %s, Starting...", path))
	}
}
