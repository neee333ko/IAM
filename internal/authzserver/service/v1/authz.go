package v1

import (
	"context"

	"github.com/ory/ladon"
)

type AuthzService interface {
	IsAllowed(ctx context.Context, request *ladon.Request) (err error)
}

type authzService struct {
	ladon ladon.Warden
}

func NewAuthzService(a Authorizer) AuthzService {
	return &authzService{
		ladon: &ladon.Ladon{
			Manager:     NewManager(a),
			AuditLogger: NewAuthzLogger(a),
		},
	}
}

func (s *authzService) IsAllowed(ctx context.Context, request *ladon.Request) (err error) {
	return s.ladon.IsAllowed(ctx, request)
}
