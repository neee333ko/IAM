package store

import (
	"context"

	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type PolicyStore interface {
	Create(c context.Context, policy *v1.Policy, options *metav1.CreateOptions) error
	Update(c context.Context, policy *v1.Policy, options *metav1.UpdateOptions) error
	Get(c context.Context, username string, options *metav1.GetOptions) (*v1.PolicyList, error)
	GetSingle(c context.Context, name string, options *metav1.GetOptions) (*v1.Policy, error)
	List(c context.Context, options *metav1.ListOptions) (*v1.PolicyList, error)
	Delete(c context.Context, name string, options *metav1.DeleteOptions) error
	DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error
}
