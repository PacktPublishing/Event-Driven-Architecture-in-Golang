package domain

import (
	"context"
)

type StoreRepository interface {
	Load(ctx context.Context, storeID string) (*Store, error)
	Save(ctx context.Context, store *Store) error
}
