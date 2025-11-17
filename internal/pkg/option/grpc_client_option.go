package option

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type GrpcClientOption struct {
	Address string `json:"address" mapstructure:"address"`
	Cert    string `json:"cert" mapstructure:"cert"`
}

func (o *GrpcClientOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("GrpcClientOption", pflag.ExitOnError)

	fs.StringVar(&o.Address, "grpcClient.address", "127.0.0.1:50051", "address of grpc service")
	fs.StringVar(&o.Cert, "grpcClient.cert", "", "cert of grpc service")

	nfs.AddFlagSet("GrpcClientOption", fs)
}

func (o *GrpcClientOption) Validate() field.ErrorList {
	return nil
}

func (o *GrpcClientOption) Complete() error {
	return nil
}
