package store

import (
	"context"

	"github.com/ory/ladon"
)

type PolicyStore interface {
	List(ctx context.Context) (*map[string][]*ladon.DefaultPolicy, error)
}
