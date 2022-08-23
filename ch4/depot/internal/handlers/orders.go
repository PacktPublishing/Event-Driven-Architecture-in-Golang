package handlers

import (
	"eda-in-golang/depot/internal/application"
	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/ddd"
)

func RegisterOrderHandlers(orderHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.ShoppingListCompleted{}, orderHandlers.OnShoppingListCompleted)
}
