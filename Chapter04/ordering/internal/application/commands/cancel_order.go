package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type CancelOrder struct {
	ID string
}

type CancelOrderHandler struct {
	orders          domain.OrderRepository
	shopping        domain.ShoppingRepository
	domainPublisher ddd.EventPublisher
}

func NewCancelOrderHandler(orders domain.OrderRepository, shopping domain.ShoppingRepository, domainPublisher ddd.EventPublisher) CancelOrderHandler {
	return CancelOrderHandler{
		orders:          orders,
		shopping:        shopping,
		domainPublisher: domainPublisher,
	}
}

func (h CancelOrderHandler) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := h.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Cancel(); err != nil {
		return err
	}

	if err = h.shopping.Cancel(ctx, order.ShoppingID); err != nil {
		return err
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
