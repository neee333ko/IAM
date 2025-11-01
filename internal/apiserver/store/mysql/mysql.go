package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/neee333ko/IAM/internal/apiserver/store"
)

type Database struct {
	db *gorm.DB
}

// func NewDatabase(db *gorm.DB) {

// }

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
