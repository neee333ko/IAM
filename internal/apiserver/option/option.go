package option

import (
	"github.com/neee333ko/IAM/internal/pkg/option"
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

type Option struct {
	InsecureOp *option.InsecureOption `json:"insecure" mapstructure:"insecure"`
	SecureOp   *option.SecureOption   `json:"secure" mapstructure:"secure"`
	RunOp      *option.RunOption      `json:"run" mapstructure:"run"`
	Feature    *option.Feature        `json:"feature" mapstructure:"feature"`
	GRPCOp     *option.GRPCOption     `json:"grpc" mapstructure:"grpc"`
	MysqlOp    *option.MysqlOption    `json:"mysql" mapstructure:"mysql"`
	RedisOp    *option.RedisOption    `json:"redis" mapstructure:"redis"`
	JwtOp      *option.JwtOption      `json:"jwt" mapstructure:"jwt"`
	LogOp      *log.Options           `json:"log" mapstructure:"log"`
}

func NewOption() *Option {
	return &Option{
		InsecureOp: new(option.InsecureOption),
		SecureOp:   new(option.SecureOption),
		RunOp:      new(option.RunOption),
		Feature:    new(option.Feature),
		GRPCOp:     new(option.GRPCOption),
		MysqlOp:    new(option.MysqlOption),
		RedisOp:    new(option.RedisOption),
		JwtOp:      new(option.JwtOption),
		LogOp:      new(log.Options),
	}
}

func (o *Option) Flags() (nfs cli.NamedFlagSets) {
	o.InsecureOp.Flags(&nfs)
	o.SecureOp.Flags(&nfs)
	o.RunOp.Flags(&nfs)
	o.Feature.Flags(&nfs)
	o.GRPCOp.Flags(&nfs)
	o.MysqlOp.Flags(&nfs)
	o.RedisOp.Flags(&nfs)
	o.JwtOp.Flags(&nfs)
	nfs.AddFlagSet("LogOption", o.LogOp.Flags())

	return
}

func (o *Option) Validate() field.ErrorList {
	var errlist field.ErrorList = make(field.ErrorList, 0)
	errlist = append(errlist, o.InsecureOp.Validate()...)
	errlist = append(errlist, o.SecureOp.Validate()...)
	errlist = append(errlist, o.RunOp.Validate()...)
	errlist = append(errlist, o.Feature.Validate()...)
	errlist = append(errlist, o.GRPCOp.Validate()...)
	errlist = append(errlist, o.MysqlOp.Validate()...)
	errlist = append(errlist, o.RedisOp.Validate()...)
	errlist = append(errlist, o.JwtOp.Validate()...)
	errlist = append(errlist, o.LogOp.Validate()...)

	return errlist
}

func (o *Option) Complete() error {
	errs := make([]error, 0)

	errs = append(errs, o.InsecureOp.Complete())
	errs = append(errs, o.SecureOp.Complete())
	errs = append(errs, o.RunOp.Complete())
	errs = append(errs, o.Feature.Complete())
	errs = append(errs, o.GRPCOp.Complete())
	errs = append(errs, o.MysqlOp.Complete())
	errs = append(errs, o.RedisOp.Complete())
	errs = append(errs, o.JwtOp.Complete())
	errs = append(errs, o.LogOp.Complete())

	if len(errs) == 0 {
		return nil
	}

	return errors.NewAggregate(errs)
}
