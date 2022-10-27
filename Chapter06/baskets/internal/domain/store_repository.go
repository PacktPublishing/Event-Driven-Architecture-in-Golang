package domain

import (
	"context"
)

type StoreRepository interface {
	Find(ctx context.Context, storeID string) (*Store, error)
}
