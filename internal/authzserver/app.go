package authzserver

import (
	"github.com/neee333ko/IAM/internal/authzserver/config"
	"github.com/neee333ko/IAM/internal/authzserver/option"
	"github.com/neee333ko/IAM/pkg/app"
	"github.com/neee333ko/log"
)

const commandDesc = `Authorization server to run ladon policies which can protecting your resources.
It is written inspired by AWS IAM policiis.

Find more iam-authz-server information at:
    https://github.com/neee333ko/IAM,

Find more ladon information at:
    https://github.com/ory/ladon`

func App(name string, basename string) *app.App {
	option := option.NewOption()

	app := app.New(name,
		basename,
		app.WithShort("authzserver"),
		app.WithDiscription(commandDesc),
		app.WithOption(option),
		app.WithDefaultPositionalArgs(),
		app.WithRunFunc(run(option)),
	)

	return app
}

func run(option *option.Option) app.RunFunc {
	return func(s string) error {
		log.Init(option.LogOp)
		defer log.Flush()

		config := config.CreateConfigFromOption(option)

		server := CreateServerFromConfig(config)

		return server.PrepareRun().Run()
	}
}
