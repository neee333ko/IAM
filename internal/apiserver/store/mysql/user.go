package mysql

import (
	"context"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/fields"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

type UserDatabase struct {
	db *gorm.DB
}

func NewUserDatabase(db *gorm.DB) store.UserStore {
	return &UserDatabase{db: db}
}

func (ud *UserDatabase) Create(c context.Context, user *v1.User, options *metav1.CreateOptions) error {
	err := ud.db.Create(user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError = new(mysql.MySQLError)
		if errors.As(err, mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = errors.WithCode(code.ErrUserAlreadyExist, err.Error())
			default:
				err = errors.WithCode(code.ErrDatabase, err.Error())
			}
		}
	}

	return err
}

func (ud *UserDatabase) Update(c context.Context, user *v1.User, options *metav1.UpdateOptions) error {
	err := ud.db.Save(user).Error
	if err != nil {
		err = errors.WithCode(code.ErrDatabase, err.Error())
	}

	return err
}

func (ud *UserDatabase) Get(c context.Context, name string, options *metav1.GetOptions) (*v1.User, error) {
	var user *v1.User = new(v1.User)
	err := ud.db.Where("name = ? and status = 1", name).First(user).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, err
}

func (ud *UserDatabase) List(c context.Context, options *metav1.ListOptions) (*v1.UserList, error) {
	fieldSelector := options.FieldSelector
	selector, err := fields.ParseSelector(fieldSelector)
	if err != nil {
		return nil, errors.WithCode(code.ErrQueryInvalid, err.Error())
	}

	value, _ := selector.RequiresExactMatch("name")

	var userlist *v1.UserList = new(v1.UserList)

	err = ud.db.Where("name = ? and status = 1", value).
		Limit(options.Limit).
		Offset(options.Offset).
		Find(userlist.Items).
		Limit(-1).Offset(-1).
		Count(userlist.TotalCount).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return userlist, nil
}

func (ud *UserDatabase) Delete(c context.Context, name string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = ud.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	err := db.Where("name = ?", name).Delete(v1.User{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	sd := &SecretDatabase{db: ud.db}
	err = sd.DeleteByUsername(c, name, options)
	if err != nil {
		return err
	}

	pd := &PolicyDatabase{db: ud.db}
	err = pd.DeleteByUsername(c, name, options)
	if err != nil {
		return err
	}

	return nil
}

func (ud *UserDatabase) DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error {
	var db *gorm.DB = ud.db
	if options.Unscoped {
		db = db.Unscoped()
	}

	var wg sync.WaitGroup

	errs := make([]error, 0)

	var mu *sync.Mutex = new(sync.Mutex)

	for _, name := range names {
		wg.Add(1)

		go func(name string) {
			defer wg.Done()

			err := db.Where("name = ?", name).Delete(v1.User{}).Error
			if err != nil {
				appendError(&errs, err, mu)
				return
			}

			sd := &SecretDatabase{db: ud.db}
			err = sd.DeleteByUsername(c, name, options)
			if err != nil {
				appendError(&errs, err, mu)
				return
			}

			pd := &PolicyDatabase{db: ud.db}
			err = pd.DeleteByUsername(c, name, options)
			if err != nil {
				appendError(&errs, err, mu)
				return
			}
		}(name)
	}

	wg.Wait()

	if len(errs) != 0 {
		return errors.WithCode(code.ErrDatabase, errors.NewAggregate(errs).Error())
	}

	return nil
}

func appendError(errs *[]error, err error, mutex *sync.Mutex) {
	mutex.Lock()
	*errs = append(*errs, err)
	mutex.Unlock()
}
