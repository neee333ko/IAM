package apiserver

import (
	"github.com/neee333ko/IAM/internal/apiserver/config"
	"github.com/neee333ko/IAM/internal/apiserver/option"
	"github.com/neee333ko/IAM/pkg/app"
	"github.com/neee333ko/log"
)

const commandDesc = ` The IAM API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.`

func App(name string, basename string) *app.App {
	option := option.NewOption()

	app := app.New(name,
		basename,
		app.WithShort("apiserver"),
		app.WithDiscription(commandDesc),
		app.WithDefaultPositionalArgs(),
		app.WithOption(option),
		app.WithRunFunc(Run(option)),
	)

	return app
}

func Run(option *option.Option) func(basename string) error {
	return func(basename string) error {
		log.Init(option.LogOp)
		log.Flush()

		config, err := config.NewConfig(option)
		if err != nil {
			return err
		}

		server, err := CreateServerFromConfig(config)
		if err != nil {
			return err
		}

		return server.PreparedRun().Run()
	}
}
