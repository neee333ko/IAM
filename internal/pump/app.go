package pump

import (
	genericserver "github.com/neee333ko/IAM/internal/pkg/server"
	"github.com/neee333ko/IAM/internal/pump/config"
	"github.com/neee333ko/IAM/internal/pump/option"
	"github.com/neee333ko/IAM/pkg/app"
	"github.com/neee333ko/log"
)

const commandDisc = `IAM Pump Server pump the record in redis to mongodb.
See document in github.com/neee333ko/IAM`

func App(name string, basename string) *app.App {
	option := option.NewOption()

	app := app.New(
		name,
		basename,
		app.WithDiscription(commandDisc),
		app.WithDefaultPositionalArgs(),
		app.WithOption(option),
		app.WithRunFunc(run(option)),
	)

	return app
}

func run(option *option.Option) app.RunFunc {
	return func(s string) error {
		log.Init(option.LogOp)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOption(option)
		if err != nil {
			return err
		}
		server := CreateServerFromConfig(cfg)

		c := genericserver.SetupSignalHandler()

		return server.PrepareRun().Run(c)
	}
}
