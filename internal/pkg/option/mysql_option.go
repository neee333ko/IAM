package option

import (
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type MysqlOption struct {
	Username        string `json:"username" mapstructure:"username"`
	Password        string `json:"password" mapstructure:"password"`
	Address         string `json:"address" mapstructure:"address"`
	DBName          string `json:"dbname" mapstructure:"dbname"`
	MaxOpenConns    int    `json:"maxOpenConns" mapstructure:"maxOpenConns"`
	MaxIdleConns    int    `json:"maxIdleConns" mapstructure:"maxIdleConns"`
	ConnMaxLifetime int64  `json:"connMaxLifeTime" mapstructure:"connMaxLifeTime"`
}

func (o *MysqlOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("MysqlOption", pflag.ExitOnError)

	fs.StringVar(&o.Username, "mysql.username", "", "username of mysql")
	fs.StringVar(&o.Password, "mysql.password", "", "password of mysql user")
	fs.StringVar(&o.Address, "mysql.address", "127.0.0.1:3306", "address of mysql service")
	fs.StringVar(&o.DBName, "mysql.dbname", "", "dbname to connect")
	fs.IntVar(&o.MaxOpenConns, "mysql.maxOpenConns", 0, "maximum number of open connections to the database")
	fs.IntVar(&o.MaxIdleConns, "mysql.maxIdleConns", 2, "maximum number of connections in the idle connection pool")
	fs.Int64Var(&o.ConnMaxLifetime, "mysql.connMaxLifeTime", 0, "maximum amount of time a connection may be reused")

	nfs.AddFlagSet("MysqlOption", fs)
}

func (o *MysqlOption) Validate() field.ErrorList {
	return nil
}

func (o *MysqlOption) Complete() error {
	if o.Address == "" {
		o.Address = "127.0.0.1:3306"
	}
	return nil
}
