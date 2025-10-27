package store

import (
	"context"

	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type UserStore interface {
	Create(c context.Context, user *v1.User, options *metav1.CreateOptions) error
	Update(c context.Context, user *v1.User, options *metav1.UpdateOptions) error
	Get(c context.Context, name string, options *metav1.GetOptions) (*v1.User, error)
	List(c context.Context, options *metav1.ListOptions) (*v1.UserList, error)
	Delete(c context.Context, name string, options *metav1.DeleteOptions) error
	DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error
}
