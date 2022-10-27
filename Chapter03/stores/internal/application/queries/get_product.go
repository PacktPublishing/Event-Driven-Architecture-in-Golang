package queries

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type GetProduct struct {
	ID string
}

type GetProductHandler struct {
	products domain.ProductRepository
}

func NewGetProductHandler(products domain.ProductRepository) GetProductHandler {
	return GetProductHandler{products: products}
}

func (h GetProductHandler) GetProduct(ctx context.Context, query GetProduct) (*domain.Product, error) {
	return h.products.FindProduct(ctx, query.ID)
}
