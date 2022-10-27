package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

func RegisterNotificationHandlers(notificationHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.OrderCreatedEvent, notificationHandlers)
	domainSubscriber.Subscribe(domain.OrderReadiedEvent, notificationHandlers)
	domainSubscriber.Subscribe(domain.OrderCanceledEvent, notificationHandlers)
}
