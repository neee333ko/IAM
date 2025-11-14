package store

import (
	"context"

	pb "github.com/neee333ko/api/proto/v1"
)

type SecretStore interface {
	List(ctx context.Context) (*map[string]*pb.SecretInfo, error)
}
