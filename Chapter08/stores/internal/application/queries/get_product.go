package queries

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type GetProduct struct {
	ID string
}

type GetProductHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductHandler(catalog domain.CatalogRepository) GetProductHandler {
	return GetProductHandler{catalog: catalog}
}

func (h GetProductHandler) GetProduct(ctx context.Context, query GetProduct) (*domain.CatalogProduct, error) {
	return h.catalog.Find(ctx, query.ID)
}
