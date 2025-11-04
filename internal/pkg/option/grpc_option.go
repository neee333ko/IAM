package option

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type GRPCOption struct {
	Address        string `json:"address" mapstructure:"address"`
	MaxRecvMsgSize int    `json:"maxRecvMsgSize" mapstructure:"maxRecvMsgSize"`
}

func (o *GRPCOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("GRPCOption", pflag.ExitOnError)

	fs.StringVar(&o.Address, "grpc.address", "127.0.0.1:50051", "address of grpc service")
	fs.IntVar(&o.MaxRecvMsgSize, "grpc.maxRecvMsgSize", 4*1024*1024, "MaxRecvMsgSize of grpc server")

	nfs.AddFlagSet("GRPCOption", fs)
}

func (o *GRPCOption) Validate() field.ErrorList {
	return nil
}

func (o *GRPCOption) Complete() error {
	if o.Address == "" {
		o.Address = "127.0.0.1:50051"
	}
	return nil
}
