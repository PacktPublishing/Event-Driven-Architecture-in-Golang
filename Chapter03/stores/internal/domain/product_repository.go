package domain

import (
	"context"
)

type ProductRepository interface {
	FindProduct(ctx context.Context, id string) (*Product, error)
	AddProduct(ctx context.Context, product *Product) error
	RemoveProduct(ctx context.Context, id string) error
	GetCatalog(ctx context.Context, storeID string) ([]*Product, error)
}
