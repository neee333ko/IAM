package option

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type Feature struct {
	EnableMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"`
	EnableProfiling bool `json:"enable-profiling" mapstructure:"enable-profiling"`
}

func (f *Feature) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("Feature", pflag.ExitOnError)

	fs.BoolVar(&f.EnableMetrics, "enable-metrics", true, "enable prometheus metrics")
	fs.BoolVar(&f.EnableProfiling, "enable-profiling", true, "enable profiling")

	nfs.AddFlagSet("Feature", fs)
}

func (o *Feature) Validate() field.ErrorList {
	return nil
}

func (o *Feature) Complete() error {
	return nil
}
