package domain

import (
	"context"
)

type FakeCatalogRepository struct {
	products map[string]*CatalogProduct
}

var _ CatalogRepository = (*FakeCatalogRepository)(nil)

func NewFakeCatalogRepository() *FakeCatalogRepository {
	return &FakeCatalogRepository{
		products: map[string]*CatalogProduct{},
	}
}

func (r *FakeCatalogRepository) AddProduct(ctx context.Context, productID, storeID, name, description, sku string, price float64) error {
	// TODO implement me
	panic("implement me")
}

func (r *FakeCatalogRepository) Rebrand(ctx context.Context, productID, name, description string) error {
	// TODO implement me
	panic("implement me")
}

func (r *FakeCatalogRepository) UpdatePrice(ctx context.Context, productID string, delta float64) error {
	// TODO implement me
	panic("implement me")
}

func (r *FakeCatalogRepository) RemoveProduct(ctx context.Context, productID string) error {
	// TODO implement me
	panic("implement me")
}

func (r *FakeCatalogRepository) Find(ctx context.Context, productID string) (*CatalogProduct, error) {
	// TODO implement me
	panic("implement me")
}

func (r *FakeCatalogRepository) GetCatalog(ctx context.Context, storeID string) ([]*CatalogProduct, error) {
	// TODO implement me
	panic("implement me")
}
