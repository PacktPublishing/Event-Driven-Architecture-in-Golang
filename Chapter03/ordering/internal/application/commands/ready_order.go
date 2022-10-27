package commands

import (
	"context"

	"eda-in-golang/ordering/internal/domain"
)

type ReadyOrder struct {
	ID string
}

type ReadyOrderHandler struct {
	orders        domain.OrderRepository
	invoices      domain.InvoiceRepository
	notifications domain.NotificationRepository
}

func NewReadyOrderHandler(orders domain.OrderRepository, invoices domain.InvoiceRepository,
	notifications domain.NotificationRepository,
) ReadyOrderHandler {
	return ReadyOrderHandler{
		orders:        orders,
		invoices:      invoices,
		notifications: notifications,
	}
}

func (h ReadyOrderHandler) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := h.orders.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = order.Ready(); err != nil {
		return nil
	}

	if err = h.orders.Update(ctx, order); err != nil {
		return err
	}

	if err = h.notifications.NotifyOrderReady(ctx, order.ID, order.CustomerID); err != nil {
		return err
	}

	return h.orders.Update(ctx, order)
}
