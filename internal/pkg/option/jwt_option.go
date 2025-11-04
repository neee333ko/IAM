package option

import (
	"time"

	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type JwtOption struct {
	Realm            string `json:"realm" mapstructure:"realm"`
	SigningAlgorithm string `json:"alg" mapstructure:"alg"`
	Key              string `json:"key" mapstructure:"key"`
	Timeout          int64  `json:"timeout" mapstructure:"timeout"`
	MaxRefresh       int64  `json:"maxRefresh" mapstructure:"maxRefresh"`
}

func (o *JwtOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("JwtOption", pflag.ExitOnError)

	fs.StringVar(&o.Realm, "jwt.realm", "", "Realm name to display to the user.")
	fs.StringVar(&o.SigningAlgorithm, "jwt.alg", "HS256", "signing algorithm - possible values are HS256, HS384, HS512, RS256, RS384 or RS512")
	fs.StringVar(&o.Key, "jwt.key", "IAMKEY", "Secret key used for signing.")
	fs.Int64Var(&o.Timeout, "jwt.timeout", int64(time.Hour), "Duration that a jwt token is valid.")
	fs.Int64Var(&o.MaxRefresh, "jwt.maxRefresh", 0, "allows clients to refresh their token until MaxRefresh has passed")

	nfs.AddFlagSet("JwtOption", fs)
}

func (o *JwtOption) Validate() field.ErrorList {
	return nil
}

func (o *JwtOption) Complete() error {
	return nil
}
