package commands

import (
	"context"

	"eda-in-golang/ordering/internal/domain"
)

type CancelOrder struct {
	ID string
}

type CancelOrderHandler struct {
	orders        domain.OrderRepository
	shopping      domain.ShoppingRepository
	notifications domain.NotificationRepository
}

func NewCancelOrderHandler(orders domain.OrderRepository, shopping domain.ShoppingRepository,
	notifications domain.NotificationRepository,
) CancelOrderHandler {
	return CancelOrderHandler{
		orders:        orders,
		shopping:      shopping,
		notifications: notifications,
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

	if err = h.notifications.NotifyOrderCanceled(ctx, order.ID, order.CustomerID); err != nil {
		return err
	}

	return h.orders.Update(ctx, order)
}
