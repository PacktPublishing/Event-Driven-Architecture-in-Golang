package handlers

import (
	"eda-in-golang/baskets/internal/application"
	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/ddd"
)

func RegisterOrderHandlers(orderHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.BasketCheckedOut{}, orderHandlers.OnBasketCheckedOut)
}
