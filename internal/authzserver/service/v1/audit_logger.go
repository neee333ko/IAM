package v1

import (
	"context"

	"github.com/ory/ladon"
)

type authzLogger struct{}

func NewAuthzLogger() ladon.AuditLogger {
	return &authzLogger{}
}

func (l *authzLogger) LogRejectedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
}

func (l *authzLogger) LogGrantedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
}
