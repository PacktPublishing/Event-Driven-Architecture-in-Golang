package commands

import (
	"context"

	"eda-in-golang/ordering/internal/domain"
)

type CancelOrder struct {
	ID string
}

type CancelOrderHandler struct {
	orders   domain.OrderRepository
	shopping domain.ShoppingRepository
}

func NewCancelOrderHandler(orders domain.OrderRepository, shopping domain.ShoppingRepository) CancelOrderHandler {
	return CancelOrderHandler{
		orders:   orders,
		shopping: shopping,
	}
}

func (h CancelOrderHandler) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Cancel(); err != nil {
		return err
	}

	if err = h.shopping.Cancel(ctx, order.ShoppingID); err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
