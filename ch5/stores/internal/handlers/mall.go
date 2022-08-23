package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

func RegisterMallHandlers(mallHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.StoreCreatedEvent, mallHandlers)
	domainSubscriber.Subscribe(domain.StoreParticipationEnabledEvent, mallHandlers)
	domainSubscriber.Subscribe(domain.StoreParticipationDisabledEvent, mallHandlers)
	domainSubscriber.Subscribe(domain.StoreRebrandedEvent, mallHandlers)
}
