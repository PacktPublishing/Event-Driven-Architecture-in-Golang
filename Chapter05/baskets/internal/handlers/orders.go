package handlers

import (
	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/ddd"
)

func RegisterOrderHandlers(orderHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(domain.BasketCheckedOutEvent, orderHandlers)
}
