package mysql

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/neee333ko/IAM/internal/apiserver/store"
	v1 "github.com/neee333ko/api/apiserver/v1"
)

type PolicyAuditDatabase struct {
	db *gorm.DB
}

func NewPolicyAuditDatabase(db *gorm.DB) store.PolicyAuditStore {
	return &PolicyAuditDatabase{db: db}
}

func (pad *PolicyAuditDatabase) ClearOutDated(ctx context.Context, maxReservedDays int) (int64, error) {
	var numbers int64
	err := pad.db.Where("createdAt < ?", time.Now().Add(-time.Duration(maxReservedDays))).Delete(v1.Policy{}).Count(&numbers).Error
	if err != nil {
		return 0, err
	}

	return numbers, nil
}
