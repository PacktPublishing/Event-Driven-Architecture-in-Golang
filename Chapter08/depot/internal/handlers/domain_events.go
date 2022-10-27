package handlers

import (
	"context"

	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/ddd"
)

type domainHandlers[T ddd.AggregateEvent] struct {
	orders domain.OrderRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*domainHandlers[ddd.AggregateEvent])(nil)

func NewDomainEventHandlers(orders domain.OrderRepository) ddd.EventHandler[ddd.AggregateEvent] {
	return domainHandlers[ddd.AggregateEvent]{
		orders: orders,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.AggregateEvent], handler ddd.EventHandler[ddd.AggregateEvent]) {
	subscriber.Subscribe(handler, domain.ShoppingListCompletedEvent)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ShoppingListCompletedEvent:
		return h.onShoppingListCompleted(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onShoppingListCompleted(ctx context.Context, event ddd.AggregateEvent) error {
	completed := event.Payload().(*domain.ShoppingListCompleted)
	return h.orders.Ready(ctx, completed.ShoppingList.OrderID)
}
