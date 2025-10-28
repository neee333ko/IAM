package v1

import (
	"context"

	"github.com/neee333ko/IAM/internal/apiserver/store"
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type UserServ interface {
	Create(c context.Context, user *v1.User, options *metav1.CreateOptions) error
	Update(c context.Context, user *v1.User, options *metav1.UpdateOptions) error
	Get(c context.Context, name string, options *metav1.GetOptions) (*v1.User, error)
	List(c context.Context, options *metav1.ListOptions) (*v1.UserList, error)
	Delete(c context.Context, name string, options *metav1.DeleteOptions) error
	DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error
}

type UserService struct {
	store store.Factory
}

func (us *UserService) Create(c context.Context, user *v1.User, options *metav1.CreateOptions) error {
	if err := us.store.NewUserStore().Create(c, user, options); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Update(c context.Context, user *v1.User, options *metav1.UpdateOptions) error {
	if err := us.store.NewUserStore().Update(c, user, options); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Get(c context.Context, name string, options *metav1.GetOptions) (*v1.User, error) {
	user, err := us.store.NewUserStore().Get(c, name, options)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) List(c context.Context, options *metav1.ListOptions) (*v1.UserList, error) {
	userlist, err := us.store.NewUserStore().List(c, options)
	if err != nil {
		return nil, err
	}

	return userlist, nil
}

func (us *UserService) Delete(c context.Context, name string, options *metav1.DeleteOptions) error {
	if err := us.store.NewUserStore().Delete(c, name, options); err != nil {
		return err
	}

	return nil
}

func (us *UserService) DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error {
	if err := us.store.NewUserStore().DeleteCollection(c, names, options); err != nil {
		return err
	}

	return nil
}
