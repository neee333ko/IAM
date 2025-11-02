package cache

import (
	"context"
	"sync"
	"time"

	"github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/pkg/code"
	pb "github.com/neee333ko/api/proto/v1"
	v1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

type CacheController struct {
	store store.Factory
}

var (
	cacheserver *CacheController
	once        sync.Once
)

func GetCacheInsOr(store store.Factory) *CacheController {
	if store != nil {
		once.Do(func() {
			cacheserver = &CacheController{store: store}
		})
	}

	if cacheserver == nil {
		log.Fatal("no cache server provided")
	}

	return cacheserver
}

func (cc *CacheController) ListSecrets(ctx context.Context, req *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	log.L(ctx).Info("grpc list secrets request received...")

	sl, err := cc.store.NewSecretStore().List(ctx, &v1.ListOptions{Offset: req.GetOffset(), Limit: req.GetLimit()})
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}

	var resp *pb.ListSecretsResponse = new(pb.ListSecretsResponse)
	resp.TotalCount = sl.GetTotalCount()

	for _, item := range sl.Items {
		resp.Items = append(resp.Items, &pb.SecretInfo{
			Name:        item.Name,
			SecretId:    item.SecretID,
			Username:    item.Username,
			SecretKey:   item.SecretKey,
			Expires:     item.Expires,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}

func (cc *CacheController) ListPolicies(ctx context.Context, req *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	log.L(ctx).Info("grpc list policies request received...")

	pl, err := cc.store.NewPolicyStore().List(ctx, &v1.ListOptions{Offset: req.GetOffset(), Limit: req.GetLimit()})
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}

	var resp *pb.ListPoliciesResponse = new(pb.ListPoliciesResponse)
	resp.TotalCount = pl.GetTotalCount()

	for _, item := range pl.Items {
		resp.Items = append(resp.Items, &pb.PolicyInfo{
			Name:         item.Name,
			Username:     item.Username,
			PolicyShadow: item.PolicyShadow,
			CreatedAt:    item.CreatedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}
