package queries

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type GetCatalog struct {
	StoreID string
}

type GetCatalogHandler struct {
	catalog domain.CatalogRepository
}

func NewGetCatalogHandler(catalog domain.CatalogRepository) GetCatalogHandler {
	return GetCatalogHandler{catalog: catalog}
}

func (h GetCatalogHandler) GetCatalog(ctx context.Context, query GetCatalog) ([]*domain.CatalogProduct, error) {
	return h.catalog.GetCatalog(ctx, query.StoreID)
}
