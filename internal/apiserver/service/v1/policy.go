package v1

import (
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type PolicyServ interface {
	Create(policy *v1.Policy, options *metav1.CreateOptions) error
	Update(policy *v1.Policy, options *metav1.UpdateOptions) error
	Get(name string, options *metav1.GetOptions) (*v1.Policy, error)
	List(options *metav1.ListOptions) (*v1.PolicyList, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(policylist *v1.PolicyList, options *metav1.DeleteOptions) error
}
