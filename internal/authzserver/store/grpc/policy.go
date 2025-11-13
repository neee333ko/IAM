package grpc

import (
	"context"

	pb "github.com/neee333ko/api/proto/v1"
	"github.com/neee333ko/component-base/pkg/json"
	"github.com/ory/ladon"
)

type PolicyClient struct {
	pb.CacheClient
}

func (pc *PolicyClient) List(ctx context.Context) (*map[string][]*ladon.DefaultPolicy, error) {
	var limit int64 = -1
	var offset int64 = -1

	resp, err := pc.ListPolicies(ctx, &pb.ListPoliciesRequest{Limit: &limit, Offset: &offset})
	if err != nil {
		return nil, err
	}

	policies := make(map[string][]*ladon.DefaultPolicy, 0)

	for _, item := range resp.Items {
		var policy ladon.DefaultPolicy

		_ = json.Unmarshal([]byte(item.PolicyShadow), &policy)

		policies[item.Username] = append(policies[item.Username], &policy)
	}

	return &policies, nil
}
