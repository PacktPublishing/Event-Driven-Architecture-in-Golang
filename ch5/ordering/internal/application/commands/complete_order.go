package commands

import (
	"context"

	"eda-in-golang/ordering/internal/domain"
)

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type CompleteOrderHandler struct {
	orders domain.OrderRepository
}

func NewCompleteOrderHandler(orders domain.OrderRepository) CompleteOrderHandler {
	return CompleteOrderHandler{
		orders: orders,
	}
}

func (h CompleteOrderHandler) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = order.Complete(cmd.InvoiceID)
	if err != nil {
		return nil
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return nil
}
