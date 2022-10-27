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
	orders    domain.OrderRepository
	shopping  domain.ShoppingRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCancelOrderHandler(orders domain.OrderRepository, shopping domain.ShoppingRepository, publisher ddd.EventPublisher[ddd.Event]) CancelOrderHandler {
	return CancelOrderHandler{
		orders:    orders,
		shopping:  shopping,
		publisher: publisher,
	}
}

func (h CancelOrderHandler) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Cancel()
	if err != nil {
		return err
	}

	// TODO CH8 remove; handled in the saga
	if err = h.shopping.Cancel(ctx, order.ShoppingID); err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
