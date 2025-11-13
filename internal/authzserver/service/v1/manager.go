package v1

import (
	"context"

	"github.com/ory/ladon"
)

type authzManager struct {
	a Authorizer
}

func NewManager(a Authorizer) ladon.Manager {
	return &authzManager{a: a}
}

func (m *authzManager) Create(ctx context.Context, policy ladon.Policy) error {
	return nil
}

// Update updates an existing policy.
func (m *authzManager) Update(ctx context.Context, policy ladon.Policy) error {
	return nil
}

// Get retrieves a policy.
func (m *authzManager) Get(ctx context.Context, id string) (ladon.Policy, error) {
	return nil, nil
}

// Delete removes a policy.
func (m *authzManager) Delete(ctx context.Context, id string) error {
	return nil
}

// GetAll retrieves all policies.
func (m *authzManager) GetAll(ctx context.Context, limit, offset int64) (ladon.Policies, error) {
	return nil, nil
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
func (m *authzManager) FindRequestCandidates(ctx context.Context, r *ladon.Request) (ladon.Policies, error) {
	return m.a.FindRequestCandidates(ctx, r)
}

// FindPoliciesForSubject returns policies that could match the subject. It either returns
// a set of policies that applies to the subject, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *authzManager) FindPoliciesForSubject(ctx context.Context, subject string) (ladon.Policies, error) {
	return nil, nil
}

// FindPoliciesForResource returns policies that could match the resource. It either returns
// a set of policies that apply to the resource, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *authzManager) FindPoliciesForResource(ctx context.Context, resource string) (ladon.Policies, error) {
	return nil, nil
}
