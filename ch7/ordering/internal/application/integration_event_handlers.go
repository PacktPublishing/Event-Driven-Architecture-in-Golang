package application

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
	"eda-in-golang/ordering/orderingpb"
)

type IntegrationEventHandlers[T ddd.AggregateEvent] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*IntegrationEventHandlers[ddd.AggregateEvent])(nil)

func NewIntegrationEventHandlers(publisher am.MessagePublisher[ddd.Event]) *IntegrationEventHandlers[ddd.AggregateEvent] {
	return &IntegrationEventHandlers[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func (h IntegrationEventHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case domain.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case domain.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	case domain.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)
	}
	return nil
}

func (h IntegrationEventHandlers[T]) onOrderCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.OrderCreated)
	items := make([]*orderingpb.OrderCreated_Item, len(payload.Items))
	for i, item := range payload.Items {
		items[i] = &orderingpb.OrderCreated_Item{
			ProductId: item.ProductID,
			Price:     item.Price,
			Quantity:  int32(item.Quantity),
		}
	}
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderCreatedEvent, &orderingpb.OrderCreated{
			Id:         event.AggregateID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
			ShoppingId: payload.ShoppingID,
			Items:      items,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onOrderReadied(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.OrderReadied)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderReadiedEvent, &orderingpb.OrderReadied{
			Id:         event.AggregateID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
			Total:      payload.Total,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onOrderCanceled(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.OrderCanceled)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderCanceledEvent, &orderingpb.OrderCanceled{
			Id:         event.AggregateID(),
			CustomerId: payload.CustomerID,
			PaymentId:  payload.PaymentID,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onOrderCompleted(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.OrderCompleted)
	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
		ddd.NewEvent(orderingpb.OrderCompletedEvent, &orderingpb.OrderCompleted{
			Id:         event.AggregateID(),
			CustomerId: payload.CustomerID,
			InvoiceId:  payload.InvoiceID,
		}),
	)
}
