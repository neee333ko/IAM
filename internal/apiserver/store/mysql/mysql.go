package mysql

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/pkg/option"
	"github.com/neee333ko/IAM/pkg/db"
)

type Database struct {
	db *gorm.DB
}

var (
	mysqlDB *Database
	once    sync.Once
)

func GetMysqlInsOr(option *option.MysqlOption) store.Factory {
	if option != nil {
		once.Do(func() {
			o := &db.Option{
				Username:        option.Username,
				Password:        option.Password,
				Address:         option.Address,
				DBName:          option.DBName,
				MaxOpenConns:    option.MaxOpenConns,
				MaxIdleConns:    option.MaxIdleConns,
				ConnMaxLifetime: time.Duration(option.ConnMaxLifetime),
			}

			db := db.NewDB(o)

			mysqlDB = &Database{db: db}
		})
	}

	return mysqlDB
}

func (database *Database) NewUserStore() store.UserStore {
	return NewUserDatabase(database.db)
}

func (database *Database) NewSecretStore() store.SecretStore {
	return NewSecretDatabase(database.db)
}

func (database *Database) NewPolicyStore() store.PolicyStore {
	return NewPolicyDatabase(database.db)
}

func (database *Database) NewPolicyAuditStore() store.PolicyAuditStore {
	return NewPolicyAuditDatabase(database.db)
}

func (database *Database) Close() error {
	return database.db.Close()
}
