package store

import "context"

type PolicyAuditStore interface {
	ClearOutDated(ctx context.Context, maxReservedDays int) (int64, error)
}
