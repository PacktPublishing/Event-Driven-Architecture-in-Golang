package domain

import (
	"context"
)

type FakeProductRepository struct {
	products map[string]*Product
}

func NewFakeProductRepository() *FakeProductRepository {
	return &FakeProductRepository{products: map[string]*Product{}}
}

var _ ProductRepository = (*FakeProductRepository)(nil)

func (r *FakeProductRepository) Load(ctx context.Context, productID string) (*Product, error) {
	if product, exists := r.products[productID]; exists {
		return product, nil
	}

	return NewProduct(productID), nil
}

func (r *FakeProductRepository) Save(ctx context.Context, product *Product) error {
	for _, event := range product.Events() {
		if err := product.ApplyEvent(event); err != nil {
			return err
		}
	}

	r.products[product.ID()] = product

	return nil
}

func (r *FakeProductRepository) Reset(products ...*Product) {
	r.products = make(map[string]*Product)

	for _, product := range products {
		r.products[product.ID()] = product
	}
}
