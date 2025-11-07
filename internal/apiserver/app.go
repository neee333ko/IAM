package apiserver

import (
	"github.com/neee333ko/IAM/internal/apiserver/config"
	"github.com/neee333ko/IAM/internal/apiserver/option"
	"github.com/neee333ko/IAM/pkg/app"
	"github.com/neee333ko/log"
)

func App(name string, basename string) *app.App {
	option := option.NewOption()

	app := app.New(name,
		basename,
		app.WithShort("apiserver"),
		app.WithDiscription("iam apiserver, provides user,secret,policy api"),
		app.WithDefaultPositionalArgs(),
		app.WithOption(option),
		app.WithRunFunc(Run(option)),
	)

	return app
}

func Run(option *option.Option) func(basename string) {
	return func(basename string) {
		log.Init(option.LogOp)
		config := config.NewConfig(option)

		server := CreateServerFromConfig(config)

		server.PreparedRun().Run(basename)
	}
}
