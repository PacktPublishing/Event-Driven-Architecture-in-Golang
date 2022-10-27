package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type CompleteOrderHandler struct {
	orders          domain.OrderRepository
	domainPublisher ddd.EventPublisher
}

func NewCompleteOrderHandler(orders domain.OrderRepository, domainPublisher ddd.EventPublisher) CompleteOrderHandler {
	return CompleteOrderHandler{
		orders:          orders,
		domainPublisher: domainPublisher,
	}
}

func (h CompleteOrderHandler) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	order, err := h.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = order.Complete(cmd.InvoiceID)
	if err != nil {
		return nil
	}

	if err = h.orders.Update(ctx, order); err != nil {
		return err
	}

	// publish domain events
	if err = h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return err
	}

	return nil
}
