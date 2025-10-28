package v1

import (
	"context"

	"github.com/neee333ko/IAM/internal/apiserver/store"
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type SecretServ interface {
	Create(c context.Context, secret *v1.Secret, options *metav1.CreateOptions) error
	// Update(c context.Context, secret *v1.Secret, options *metav1.UpdateOptions) error
	Get(c context.Context, username string, options *metav1.GetOptions) (*v1.SecretList, error)
	GetSingle(c context.Context, name string, options *metav1.GetOptions) (*v1.Secret, error)
	List(c context.Context, options *metav1.ListOptions) (*v1.SecretList, error)
	Delete(c context.Context, name string, options *metav1.DeleteOptions) error
	DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error
}

type SecretService struct {
	store store.Factory
}

func (ss *SecretService) Create(c context.Context, secret *v1.Secret, options *metav1.CreateOptions) error {
	if err := ss.store.NewSecretStore().Create(c, secret, options); err != nil {
		return err
	}

	return nil
}

func (ss *SecretService) Get(c context.Context, username string, options *metav1.GetOptions) (*v1.SecretList, error) {
	sl, err := ss.store.NewSecretStore().Get(c, username, options)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (ss *SecretService) GetSingle(c context.Context, username string, options *metav1.GetOptions) (*v1.Secret, error) {
	s, err := ss.store.NewSecretStore().GetSingle(c, username, options)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (ss *SecretService) List(c context.Context, options *metav1.ListOptions) (*v1.SecretList, error) {
	sl, err := ss.store.NewSecretStore().List(c, options)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (ss *SecretService) Delete(c context.Context, name string, options *metav1.DeleteOptions) error {
	if err := ss.store.NewSecretStore().Delete(c, name, options); err != nil {
		return err
	}

	return nil
}

func (ss *SecretService) DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error {
	if err := ss.store.NewSecretStore().DeleteCollection(c, names, options); err != nil {
		return err
	}

	return nil
}
