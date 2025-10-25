package v1

import (
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type SecretServ interface {
	Create(secret *v1.Secret, options *metav1.CreateOptions) error
	Update(secret *v1.Secret, options *metav1.UpdateOptions) error
	Get(name string, options *metav1.GetOptions) (*v1.Secret, error)
	List(options *metav1.ListOptions) (*v1.SecretList, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(secretlist *v1.SecretList, options *metav1.DeleteOptions) error
}
