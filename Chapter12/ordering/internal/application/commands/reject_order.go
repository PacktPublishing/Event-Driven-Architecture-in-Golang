package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type RejectOrder struct {
	ID string
}

type RejectOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRejectOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) RejectOrderHandler {
	return RejectOrderHandler{
		orders:    orders,
		publisher: publisher,
	}
}

func (h RejectOrderHandler) RejectOrder(ctx context.Context, cmd RejectOrder) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Reject()
	if err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
