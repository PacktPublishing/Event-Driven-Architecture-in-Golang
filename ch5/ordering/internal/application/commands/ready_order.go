package commands

import (
	"context"

	"eda-in-golang/ordering/internal/domain"
)

type ReadyOrder struct {
	ID string
}

type ReadyOrderHandler struct {
	orders domain.OrderRepository
}

func NewReadyOrderHandler(orders domain.OrderRepository) ReadyOrderHandler {
	return ReadyOrderHandler{
		orders: orders,
	}
}

func (h ReadyOrderHandler) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Ready(); err != nil {
		return nil
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
