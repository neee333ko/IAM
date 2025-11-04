package option

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type InsecureOption struct {
	Address string `json:"address" mapstructure:"address"`
}

func (o *InsecureOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("InsecureOption", pflag.ExitOnError)

	fs.StringVar(&o.Address, "insecure.address", "127.0.0.1:80", "address of http service")

	nfs.AddFlagSet("InsecureOption", fs)
}

func (o *InsecureOption) Validate() field.ErrorList {
	return nil
}

func (o *InsecureOption) Complete() error {
	if o.Address == "" {
		o.Address = "127.0.0.1:80"
	}

	return nil
}
