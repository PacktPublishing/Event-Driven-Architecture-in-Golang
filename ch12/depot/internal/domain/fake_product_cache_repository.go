package domain

import (
	"context"

	"github.com/stackus/errors"
)

type FakeProductCacheRepository struct {
	products map[string]*Product
}

var _ ProductCacheRepository = (*FakeProductCacheRepository)(nil)

func NewFakeProductCacheRepository() *FakeProductCacheRepository {
	return &FakeProductCacheRepository{products: map[string]*Product{}}
}

func (r *FakeProductCacheRepository) Add(ctx context.Context, productID, storeID, name string) error {
	r.products[productID] = &Product{
		ID:      productID,
		StoreID: storeID,
		Name:    name,
	}

	return nil
}

func (r *FakeProductCacheRepository) Rebrand(ctx context.Context, productID, name string) error {
	if product, exists := r.products[productID]; exists {
		product.Name = name
	}

	return nil
}

func (r *FakeProductCacheRepository) Remove(ctx context.Context, productID string) error {
	delete(r.products, productID)

	return nil
}

func (r *FakeProductCacheRepository) Find(ctx context.Context, productID string) (*Product, error) {
	if product, exists := r.products[productID]; exists {
		return product, nil
	}

	return nil, errors.ErrNotFound.Msgf("product with id: `%s` does not exist", productID)
}

func (r *FakeProductCacheRepository) Reset(products ...*Product) {
	r.products = make(map[string]*Product)

	for _, product := range products {
		r.products[product.ID] = product
	}
}
