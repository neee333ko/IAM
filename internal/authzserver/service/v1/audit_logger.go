package v1

import (
	"context"

	"github.com/ory/ladon"
)

type authzLogger struct {
	a Authorizer
}

func NewAuthzLogger(a Authorizer) ladon.AuditLogger {
	return &authzLogger{a: a}
}

func (l *authzLogger) LogRejectedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	l.a.LogRejectedAccessRequest(ctx, request, pool, deciders)
}

func (l *authzLogger) LogGrantedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	l.a.LogGrantedAccessRequest(ctx, request, pool, deciders)
}
