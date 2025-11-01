package mysql

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/fields"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

type PolicyDatabase struct {
	db *gorm.DB
}

func NewPolicyDatabase(db *gorm.DB) store.PolicyStore {
	return &PolicyDatabase{db: db}
}

func (pd *PolicyDatabase) Create(c context.Context, policy *v1.Policy, options *metav1.CreateOptions) error {
	err := pd.db.Create(policy).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError = new(mysql.MySQLError)
		if errors.As(err, mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = errors.WithCode(code.ErrPolicyAlreadyExist, err.Error())
			default:
				err = errors.WithCode(code.ErrDatabase, err.Error())

			}
		} else {
			err = errors.WithCode(code.ErrUnknown, err.Error())
		}
	}

	return err
}

func (pd *PolicyDatabase) Update(c context.Context, policy *v1.Policy, options *metav1.UpdateOptions) error {
	err := pd.db.Save(policy).Error
	if err != nil {
		err = errors.WithCode(code.ErrDatabase, err.Error())
	}

	return err
}

func (pd *PolicyDatabase) Get(c context.Context, username string, options *metav1.GetOptions) (*v1.PolicyList, error) {
	var pl *v1.PolicyList = new(v1.PolicyList)
	err := pd.db.Where("username = ?", username).Find(pl.Items).Count(pl.ListMeta).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return pl, nil
}

func (pd *PolicyDatabase) GetSingle(c context.Context, name string, options *metav1.GetOptions) (*v1.Policy, error) {
	var policy *v1.Policy = new(v1.Policy)
	err := pd.db.Where("name = ?", name).First(policy).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.WithCode(code.ErrPolicyNotFound, err.Error())
		} else {
			err = errors.WithCode(code.ErrDatabase, err.Error())
		}
	}

	return policy, err
}

func (pd *PolicyDatabase) List(c context.Context, options *metav1.ListOptions) (*v1.PolicyList, error) {
	fieldSelector := options.FieldSelector
	selector, err := fields.ParseSelector(fieldSelector)
	if err != nil {
		return nil, errors.WithCode(code.ErrQueryInvalid, err.Error())
	}

	value, _ := selector.RequiresExactMatch("name")

	var pl *v1.PolicyList = new(v1.PolicyList)

	err = pd.db.Where("name = ?", value).
		Limit(options.Limit).
		Offset(options.Offset).
		Find(pl.Items).
		Limit(-1).
		Offset(-1).
		Count(pl.TotalCount).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return pl, nil
}

func (pd *PolicyDatabase) Delete(c context.Context, name string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = pd.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	err := db.Where("name = ?", name).Delete(v1.Policy{}).Error
	if err != nil {
		err = errors.WithCode(code.ErrDatabase, err.Error())
	}

	return err
}

func (pd *PolicyDatabase) DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = pd.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	err := db.Where("name in ?", names).Delete(v1.Policy{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return err
}

func (pd *PolicyDatabase) DeleteByUsername(c context.Context, username string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = pd.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	err := db.Where("username = ?", username).Delete(v1.Policy{}).Error
	if err != nil {
		err = errors.WithCode(code.ErrDatabase, err.Error())
	}

	return err
}
