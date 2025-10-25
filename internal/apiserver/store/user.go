package store

import (
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type UserStore interface {
	Create(user *v1.User, options *metav1.CreateOptions) error
	Update(user *v1.User, options *metav1.UpdateOptions) error
	Get(name string, options *metav1.GetOptions) (*v1.User, error)
	List(options *metav1.ListOptions) (*v1.UserList, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(userlist *v1.UserList, options *metav1.DeleteOptions) error
}
