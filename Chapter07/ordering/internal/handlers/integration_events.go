package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

func RegisterIntegrationEventHandlers(eventHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(eventHandlers,
		domain.OrderCreatedEvent,
		domain.OrderReadiedEvent,
		domain.OrderCanceledEvent,
		domain.OrderCompletedEvent,
	)
}
