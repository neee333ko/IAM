package store

import (
	"context"
)

type SecretStore interface {
	List(ctx context.Context) (*map[string]string, error)
}
