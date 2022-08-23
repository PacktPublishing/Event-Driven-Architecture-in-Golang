package domain

import (
	"context"
)

type ProductCacheRepository interface {
	Add(ctx context.Context, productID, storeID, name string, price float64) error
	Rebrand(ctx context.Context, productID, name string) error
	UpdatePrice(ctx context.Context, productID string, delta float64) error
	Remove(ctx context.Context, productID string) error
	ProductRepository
}
