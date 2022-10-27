package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/ordering/internal/domain"
)

type CreateOrder struct {
	ID         string
	CustomerID string
	PaymentID  string
	Items      []*domain.Item
}

type CreateOrderHandler struct {
	orders        domain.OrderRepository
	customers     domain.CustomerRepository
	payments      domain.PaymentRepository
	shopping      domain.ShoppingRepository
	notifications domain.NotificationRepository
}

func NewCreateOrderHandler(orders domain.OrderRepository, customers domain.CustomerRepository,
	payments domain.PaymentRepository, shopping domain.ShoppingRepository, notifications domain.NotificationRepository,
) CreateOrderHandler {
	return CreateOrderHandler{
		orders:        orders,
		customers:     customers,
		payments:      payments,
		shopping:      shopping,
		notifications: notifications,
	}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	order, err := domain.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return errors.Wrap(err, "create order command")
	}

	// authorizeCustomer
	if err = h.customers.Authorize(ctx, order.CustomerID); err != nil {
		return errors.Wrap(err, "order customer authorization")
	}

	// validatePayment
	if err = h.payments.Confirm(ctx, order.PaymentID); err != nil {
		return errors.Wrap(err, "order payment confirmation")
	}

	// scheduleShopping
	if order.ShoppingID, err = h.shopping.Create(ctx, order); err != nil {
		return errors.Wrap(err, "order shopping scheduling")
	}

	// notifyOrderCreated
	if err = h.notifications.NotifyOrderCreated(ctx, order.ID, order.CustomerID); err != nil {
		return err
	}

	return errors.Wrap(h.orders.Save(ctx, order), "create order command")
}
