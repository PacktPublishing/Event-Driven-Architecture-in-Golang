package queries

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/ordering/internal/domain"
)

type GetOrder struct {
	ID string
}

type GetOrderHandler struct {
	repo domain.OrderRepository
}

func NewGetOrderHandler(repo domain.OrderRepository) GetOrderHandler {
	return GetOrderHandler{repo: repo}
}

func (h GetOrderHandler) GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := h.repo.Find(ctx, query.ID)

	return order, errors.Wrap(err, "get order query")
}
