package queries

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type GetParticipatingStores struct{}

type GetParticipatingStoresHandler struct {
	mall domain.MallRepository
}

func NewGetParticipatingStoresHandler(mall domain.MallRepository) GetParticipatingStoresHandler {
	return GetParticipatingStoresHandler{mall: mall}
}

func (h GetParticipatingStoresHandler) GetParticipatingStores(ctx context.Context, _ GetParticipatingStores) ([]*domain.MallStore, error) {
	return h.mall.AllParticipating(ctx)
}
