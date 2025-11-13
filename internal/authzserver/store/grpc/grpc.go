package grpc

import (
	"sync"

	"github.com/neee333ko/IAM/internal/authzserver/store"
	pb "github.com/neee333ko/api/proto/v1"
	"github.com/neee333ko/log"
)

type GRPCClient struct {
	pb.CacheClient
}

var (
	grpcClient *GRPCClient
	once       sync.Once
)

func (c *GRPCClient) NewSecretStore() store.SecretStore {
	return &SecretClient{CacheClient: c.CacheClient}
}

func (c *GRPCClient) NewPolicyStore() store.PolicyStore {
	return &PolicyClient{CacheClient: c.CacheClient}
}

func GetGrpcClientInsOrDie(cache pb.CacheClient) store.Factory {
	if cache != nil {
		once.Do(func() {
			grpcClient = &GRPCClient{
				CacheClient: cache,
			}
		})
		return grpcClient
	}

	if grpcClient == nil {
		log.Fatal("failed to get grpc client.")
	}

	return grpcClient
}
