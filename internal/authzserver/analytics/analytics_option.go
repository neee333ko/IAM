package analytics

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type AnalyticOption struct {
	PoolSize           int64 `json:"pool-size" mapstructure:"pool-size"`
	BufferSize         int64 `json:"buffer-size" mapstructure:"buffer-size"`
	ForcedPushInterval int64 `json:"interval" mapstructure:"interval"`
}

func (a *AnalyticOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("AnalyticOption", pflag.ExitOnError)

	fs.Int64Var(&analytics.poolSize, "pool-size", 8, "total worker number in pool")
	fs.Int64Var(&a.BufferSize, "buffer-size", 20, "buffer size of each worker")
	fs.Int64Var(&a.ForcedPushInterval, "interval", 50, "max miliseconds between two push")

	nfs.AddFlagSet("AnalyticOption", fs)
}

func (a *AnalyticOption) Validate() field.ErrorList {
	return nil
}

func (a *AnalyticOption) Complete() error {
	return nil
}
