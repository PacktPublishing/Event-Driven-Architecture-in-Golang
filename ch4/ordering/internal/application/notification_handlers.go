package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type NotificationHandlers struct {
	notifications domain.NotificationRepository
	ignoreUnimplementedDomainEvents
}

var _ DomainEventHandlers = (*NotificationHandlers)(nil)

func NewNotificationHandlers(notifications domain.NotificationRepository) *NotificationHandlers {
	return &NotificationHandlers{
		notifications: notifications,
	}
}

func (h NotificationHandlers) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	orderCreated := event.(*domain.OrderCreated)
	return h.notifications.NotifyOrderCreated(ctx, orderCreated.Order.ID, orderCreated.Order.CustomerID)
}

func (h NotificationHandlers) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	orderReadied := event.(*domain.OrderReadied)
	return h.notifications.NotifyOrderReady(ctx, orderReadied.Order.ID, orderReadied.Order.CustomerID)
}

func (h NotificationHandlers) OnOrderCanceled(ctx context.Context, event ddd.Event) error {
	orderCanceled := event.(*domain.OrderCanceled)
	return h.notifications.NotifyOrderCanceled(ctx, orderCanceled.Order.ID, orderCanceled.Order.CustomerID)
}
