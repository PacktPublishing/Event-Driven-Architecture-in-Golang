package domain

import (
	"context"
)

type ProductCacheRepository interface {
	Add(ctx context.Context, productID, storeID, name string) error
	Rebrand(ctx context.Context, productID, name string) error
	Remove(ctx context.Context, productID string) error
	ProductRepository
}
