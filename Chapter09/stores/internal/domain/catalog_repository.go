package domain

import (
	"context"
)

type CatalogProduct struct {
	ID          string
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

type CatalogRepository interface {
	AddProduct(ctx context.Context, productID, storeID, name, description, sku string, price float64) error
	Rebrand(ctx context.Context, productID, name, description string) error
	UpdatePrice(ctx context.Context, productID string, delta float64) error
	RemoveProduct(ctx context.Context, productID string) error
	Find(ctx context.Context, productID string) (*CatalogProduct, error)
	GetCatalog(ctx context.Context, storeID string) ([]*CatalogProduct, error)
}
