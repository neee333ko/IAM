package option

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type RedisOption struct {
	Host                  string   `json:"host" mapstructure:"host"`
	Port                  int      `json:"port" mapstructure:"port"`
	Addrs                 []string `json:"addrs" mapstructure:"addrs"`
	MasterName            string   `json:"mastername" mapstructure:"mastername"`
	Username              string   `json:"username" mapstructure:"username"`
	Password              string   `json:"password" mapstructure:"password"`
	Database              int      `json:"database" mapstructure:"password"`
	MaxIdle               int      `json:"maxIdle" mapstructure:"maxIdle"`
	MaxActive             int      `json:"maxActive" mapstructure:"maxActive"`
	Timeout               int      `json:"timeout" mapstructure:"timeout"`
	EnableCluster         bool     `json:"enableCluster" mapstructure:"enableCluster"`
	UseSSL                bool     `json:"SSL" mapstructure:"SSL"`
	SSLInsecureSkipVerify bool     `json:"skipVerify" mapstructure:"skipVerify"`
}

func (o *RedisOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("RedisOption", pflag.ExitOnError)

	fs.StringVar(&o.Host, "redis.host", "127.0.0.1", "host of redis service")
	fs.IntVar(&o.Port, "redis.port", 6379, "port of redis service")
	fs.StringSliceVar(&o.Addrs, "redis.addrs", nil, "either a single address or a seed list of host:port addresses of cluster/sentinel nodes")
	fs.StringVar(&o.MasterName, "redis.mastername", "", "MasterName is the sentinel master name. Only for failover clients.")
	fs.StringVar(&o.Username, "redis.username", "", "username of redis")
	fs.StringVar(&o.Password, "redis.password", "", "password of redis")
	fs.IntVar(&o.Database, "redis.database", 0, "Database to be selected after connecting to the server.Only single-node and failover clients.")
	fs.IntVar(&o.MaxIdle, "redis.maxIdle", 2, "maxIdle Conns")
	fs.IntVar(&o.MaxActive, "redis.maxActive", 10, "maxActive Conns")
	fs.IntVar(&o.Timeout, "redis.timeout", 0, "timeout of client")
	fs.BoolVar(&o.EnableCluster, "redis.enableCluster", false, "enable cluster client")
	fs.BoolVar(&o.UseSSL, "redis.SSL", false, "enable SSL")
	fs.BoolVar(&o.SSLInsecureSkipVerify, "redis.skipVerify", true, "skip SSL verify")

	nfs.AddFlagSet("RedisOption", fs)
}

func (o *RedisOption) Validate() field.ErrorList {
	return nil
}

func (o *RedisOption) Complete() error {
	if o.Host == "" {
		o.Host = "127.0.0.1"
	}

	if o.Port <= 0 || o.Port > 65535 {
		o.Port = 6379
	}

	return nil
}
