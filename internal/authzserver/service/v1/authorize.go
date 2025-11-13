package v1

import (
	"context"

	"github.com/ory/ladon"
)

type Authorizer interface {
	FindRequestCandidates(ctx context.Context, r *ladon.Request) (ladon.Policies, error)

	LogRejectedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies)
	LogGrantedAccessRequest(ctx context.Context, request *ladon.Request, pool ladon.Policies, deciders ladon.Policies)
}
