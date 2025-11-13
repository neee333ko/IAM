package grpc

import (
	"context"

	pb "github.com/neee333ko/api/proto/v1"
)

type SecretClient struct {
	pb.CacheClient
}

func (sc *SecretClient) List(ctx context.Context) (*map[string]string, error) {
	var offset int64 = -1
	var limit int64 = -1
	resp, err := sc.ListSecrets(ctx, &pb.ListSecretsRequest{Offset: &offset, Limit: &limit})
	if err != nil {
		return nil, err
	}

	secrets := make(map[string]string, 0)

	for _, item := range resp.Items {
		secrets[item.SecretId] = item.SecretKey
	}

	return &secrets, nil
}
