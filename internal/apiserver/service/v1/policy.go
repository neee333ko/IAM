package v1

import (
	"context"

	"github.com/neee333ko/IAM/internal/apiserver/store"
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type PolicyServ interface {
	Create(c context.Context, policy *v1.Policy, options *metav1.CreateOptions) error
	Update(c context.Context, policy *v1.Policy, options *metav1.UpdateOptions) error
	Get(c context.Context, username string, options *metav1.GetOptions) (*v1.PolicyList, error)
	GetSingle(c context.Context, name string, options *metav1.GetOptions) (*v1.Policy, error)
	List(c context.Context, options *metav1.ListOptions) (*v1.PolicyList, error)
	Delete(c context.Context, name string, options *metav1.DeleteOptions) error
	DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error
}

type PolicyService struct {
	store store.Factory
}

func (ps *PolicyService) Create(c context.Context, policy *v1.Policy, options *metav1.CreateOptions) error {
	if err := ps.store.NewPolicyStore().Create(c, policy, options); err != nil {
		return err
	}

	return nil
}

func (ps *PolicyService) Update(c context.Context, policy *v1.Policy, options *metav1.UpdateOptions) error {
	if err := ps.store.NewPolicyStore().Update(c, policy, options); err != nil {
		return err
	}

	return nil
}

func (ps *PolicyService) Get(c context.Context, username string, options *metav1.GetOptions) (*v1.PolicyList, error) {
	pl, err := ps.store.NewPolicyStore().Get(c, username, options)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

func (ps *PolicyService) GetSingle(c context.Context, name string, options *metav1.GetOptions) (*v1.Policy, error) {
	p, err := ps.store.NewPolicyStore().GetSingle(c, name, options)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (ps *PolicyService) List(c context.Context, options *metav1.ListOptions) (*v1.PolicyList, error) {
	pl, err := ps.store.NewPolicyStore().List(c, options)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

func (ps *PolicyService) Delete(c context.Context, name string, options *metav1.DeleteOptions) error {
	if err := ps.store.NewPolicyStore().Delete(c, name, options); err != nil {
		return err
	}

	return nil
}

func (ps *PolicyService) DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error {
	if err := ps.store.NewPolicyStore().DeleteCollection(c, names, options); err != nil {
		return err
	}

	return nil
}
