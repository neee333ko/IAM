package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type Option struct {
	Username        string
	Password        string
	Address         string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewDB(op *Option) *gorm.DB {
	DSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", op.Username, op.Password, op.Address, op.DBName)

	db, err := gorm.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err.Error())
	}

	db.DB().SetConnMaxLifetime(op.ConnMaxLifetime)
	db.DB().SetMaxIdleConns(op.MaxIdleConns)
	db.DB().SetMaxOpenConns(op.MaxOpenConns)

	RegisterCallBacks(db)

	return db
}
