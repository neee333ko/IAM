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

type SecretDatabase struct {
	db *gorm.DB
}

func NewSecretDatabase(db *gorm.DB) store.SecretStore {
	return &SecretDatabase{db: db}
}

func (sd *SecretDatabase) Create(c context.Context, secret *v1.Secret, options *metav1.CreateOptions) error {
	err := sd.db.Create(secret).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError = new(mysql.MySQLError)
		if errors.As(err, mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = errors.WithCode(code.ErrSecretAlreadyExist, err.Error())
			default:
				err = errors.WithCode(code.ErrDatabase, err.Error())
			}
		} else {
			err = errors.WithCode(code.ErrUnknown, err.Error())
		}
	}

	return err
}

func (sd *SecretDatabase) Get(c context.Context, username string, options *metav1.GetOptions) (*v1.SecretList, error) {
	var sl *v1.SecretList = new(v1.SecretList)

	err := sd.db.Where("username = ?", username).Find(sl.Items).Count(sl.TotalCount).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return sl, nil
}

func (sd *SecretDatabase) GetSingle(c context.Context, name string, options *metav1.GetOptions) (*v1.Secret, error) {
	var s *v1.Secret = new(v1.Secret)

	err := sd.db.Where("name = ?", name).First(s).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.WithCode(code.ErrSecretNotFound, err.Error())
		} else {
			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}
	}

	return s, nil
}

func (sd *SecretDatabase) List(c context.Context, options *metav1.ListOptions) (*v1.SecretList, error) {
	fieldSelector := options.FieldSelector
	selector, err := fields.ParseSelector(fieldSelector)
	if err != nil {
		return nil, errors.WithCode(code.ErrQueryInvalid, err.Error())
	}

	value, _ := selector.RequiresExactMatch("name")

	var sl *v1.SecretList = new(v1.SecretList)

	err = sd.db.Where("name = ?", value).
		Limit(options.Limit).
		Offset(options.Offset).
		Find(sl.Items).
		Limit(-1).Offset(-1).
		Count(sl.TotalCount).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return sl, nil
}

func (sd *SecretDatabase) Delete(c context.Context, name string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = sd.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	err := db.Where("name = ?", name).Delete(v1.Secret{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (sd *SecretDatabase) DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = sd.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	err := db.Where("name in ?", names).Delete(v1.Secret{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (sd *SecretDatabase) DeleteByUsername(c context.Context, username string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = sd.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	err := db.Where("username = ?", username).Delete(v1.Secret{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return err
}
