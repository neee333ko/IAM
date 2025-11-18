package option

import (
	"github.com/neee333ko/IAM/internal/pkg/option"
	"github.com/neee333ko/IAM/internal/pump/analytics"
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/neee333ko/log"
	"github.com/spf13/pflag"
)

type PumpOption struct {
	Type             string                      `json:"type" mapstructure:"type"`
	Filters          *analytics.AnalyticsFilters `json:"filters" mapstructure:"filters"`
	Timeout          int64                       `json:"timeout" mapstructure:"timeout"`
	OmitDetailedPump bool                        `json:"omitDetail" mapstructure:"omitDetail"`
	Meta             map[string]interface{}      `json:"meta" mapstructure:"meta"`
}

type Option struct {
	PurgeDelay       int64                  `json:"purgeDelay" mapstructure:"purgeDelay"`
	HealthCheckAddr  string                 `json:"health" mapstructure:"health"`
	OmitDetailedPump bool                   `json:"omitDetail" mapstructure:"omitDetail"`
	PumpOp           map[string]*PumpOption `json:"pump" mapstructure:"pump"`
	RedisOp          *option.RedisOption    `json:"redis" mapstructure:"redis"`
	LogOp            *log.Options           `json:"log" mapstructure:"log"`
}

func NewOption() *Option {
	return &Option{
		PurgeDelay:       5,
		HealthCheckAddr:  "0.0.0.0:7070",
		OmitDetailedPump: false,
		PumpOp:           map[string]*PumpOption{},
		RedisOp:          &option.RedisOption{},
		LogOp:            log.InitOptions(),
	}
}

func (o *Option) Flags() (nfs cli.NamedFlagSets) {
	var fs *pflag.FlagSet = new(pflag.FlagSet)

	fs.Int64Var(&o.PurgeDelay, "purgeDelay", 5, "purge delay for pump server")
	fs.StringVar(&o.HealthCheckAddr, "health", "0.0.0.0:7070", "health check address of pump server")
	fs.BoolVar(&o.OmitDetailedPump, "omitDetail", false, "if true, the pool and deciders of the record pumped will be nil")
	nfs.AddFlagSet("PumpServerOption", fs)

	o.RedisOp.Flags(&nfs)
	logfs := o.LogOp.Flags()

	nfs.AddFlagSet("LogOption", logfs)

	return
}

func (o *Option) Validate() field.ErrorList {
	var errlist field.ErrorList = make(field.ErrorList, 0)

	errlist = append(errlist, o.RedisOp.Validate()...)
	errlist = append(errlist, o.LogOp.Validate()...)

	return errlist
}

func (o *Option) Complete() error {
	if err := o.RedisOp.Complete(); err != nil {
		return err
	}

	if err := o.LogOp.Complete(); err != nil {
		return err
	}

	return nil
}
