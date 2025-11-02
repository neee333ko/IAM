package db

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/neee333ko/log"
)

var timekey string = "start_time"

func RegisterCallBacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("plugin:run_before_create", before)
	db.Callback().Update().Before("gorm:update").Register("plugin:run_before_update", before)
	db.Callback().Query().Before("gorm:query").Register("plugin:run_before_query", before)
	db.Callback().Delete().Before("gorm:delete").Register("plugin:run_before_delete", before)
	db.Callback().RowQuery().Before("gorm:row").Register("plugin:run_before_row", before)

	db.Callback().Create().After("gorm:create").Register("plugin:run_after_create", after)
	db.Callback().Update().After("gorm:update").Register("plugin:run_after_update", after)
	db.Callback().Query().After("gorm:query").Register("plugin:run_after_query", after)
	db.Callback().Delete().After("gorm:delete").Register("plugin:run_after_delete", after)
	db.Callback().RowQuery().After("gorm:row").Register("plugin:run_after_row", after)
}

func before(scope *gorm.Scope) {
	scope.InstanceSet(timekey, time.Now())
}

func after(scope *gorm.Scope) {
	v, ok := scope.InstanceGet(timekey)
	if !ok {
		scope.Err(errors.New("missing starting time"))
	}

	st, _ := v.(time.Time)
	et := time.Now()

	processed_time := et.Sub(st)

	log.Infof("sql processed time: %v\n", processed_time.Seconds())
}
