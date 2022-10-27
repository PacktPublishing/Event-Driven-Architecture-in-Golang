package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type ApproveOrder struct {
	ID         string
	ShoppingID string
}

type ApproveOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewApproveOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) ApproveOrderHandler {
	return ApproveOrderHandler{
		orders:    orders,
		publisher: publisher,
	}
}

func (h ApproveOrderHandler) ApproveOrder(ctx context.Context, cmd ApproveOrder) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := order.Approve(cmd.ShoppingID)
	if err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
