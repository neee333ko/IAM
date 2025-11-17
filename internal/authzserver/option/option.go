package option

import (
	"github.com/neee333ko/IAM/internal/authzserver/analytics"
	"github.com/neee333ko/IAM/internal/pkg/option"
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/neee333ko/log"
)

type Option struct {
	InsecureOp   *option.InsecureOption    `json:"insecure" mapstructure:"insecure"`
	SecureOp     *option.SecureOption      `json:"secure" mapstructure:"secure"`
	RunOp        *option.RunOption         `json:"run" mapstructure:"run"`
	Feature      *option.Feature           `json:"featrue" mapstructure:"feature"`
	GrpcClientOp *option.GrpcClientOption  `json:"grpcClient" mapstructure:"grpcClient"`
	RedisOp      *option.RedisOption       `json:"redis" mapstructure:"redis"`
	AnalyticsOp  *analytics.AnalyticOption `json:"analytics" mapstructure:"analytics"`
	LogOp        *log.Options              `json:"log" mapstructure:"log"`
}

func NewOption() *Option {
	return &Option{
		InsecureOp:   new(option.InsecureOption),
		SecureOp:     new(option.SecureOption),
		RunOp:        new(option.RunOption),
		Feature:      new(option.Feature),
		GrpcClientOp: new(option.GrpcClientOption),
		RedisOp:      new(option.RedisOption),
		AnalyticsOp:  new(analytics.AnalyticOption),
		LogOp:        log.InitOptions(),
	}
}

func (o *Option) Flags() (nfs cli.NamedFlagSets) {
	o.InsecureOp.Flags(&nfs)
	o.SecureOp.Flags(&nfs)
	o.RunOp.Flags(&nfs)
	o.Feature.Flags(&nfs)
	o.GrpcClientOp.Flags(&nfs)
	o.RedisOp.Flags(&nfs)
	o.AnalyticsOp.Flags(&nfs)
	nfs.AddFlagSet("LogOption", o.LogOp.Flags())

	return
}

func (o *Option) Validate() field.ErrorList {
	var errlist field.ErrorList

	errlist = append(errlist, o.InsecureOp.Validate()...)
	errlist = append(errlist, o.SecureOp.Validate()...)
	errlist = append(errlist, o.RunOp.Validate()...)
	errlist = append(errlist, o.Feature.Validate()...)
	errlist = append(errlist, o.GrpcClientOp.Validate()...)
	errlist = append(errlist, o.RedisOp.Validate()...)
	errlist = append(errlist, o.AnalyticsOp.Validate()...)

	return errlist
}

func (o *Option) Complete() error {
	if err := o.InsecureOp.Complete(); err != nil {
		return err
	}
	if err := o.SecureOp.Complete(); err != nil {
		return err
	}
	if err := o.RunOp.Complete(); err != nil {
		return err
	}
	if err := o.Feature.Complete(); err != nil {
		return err
	}
	if err := o.GrpcClientOp.Complete(); err != nil {
		return err
	}
	if err := o.RedisOp.Complete(); err != nil {
		return err
	}
	if err := o.AnalyticsOp.Complete(); err != nil {
		return err
	}

	return nil
}
