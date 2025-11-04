package option

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type RunOption struct {
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
	TimeOut     int64    `json:"timeout" mapstructure:"timeout"`
	Mode        string   `json:"mode" mapstructure:"mode"`
	Healthz     bool     `json:"healthz" mapstructure:"healthz"`
	Version     bool     `json:"version" mapstructure:"version"`
}

func (o *RunOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("RunOption", pflag.ExitOnError)

	fs.StringSliceVar(&o.Middlewares, "run.middlewares", nil, "middlewares used by api server")
	fs.Int64Var(&o.TimeOut, "run.timeout", 0, "timeout for closing server")
	fs.StringVar(&o.Mode, "run.mode", "release", "gin mode")
	fs.BoolVar(&o.Healthz, "run.healthz", true, "turn on health api")
	fs.BoolVar(&o.Version, "run.version", true, "turn on version api")

	nfs.AddFlagSet("RunOption", fs)
}

func (o *RunOption) Validate() field.ErrorList {
	return nil
}

func (o *RunOption) Complete() error {
	return nil
}
