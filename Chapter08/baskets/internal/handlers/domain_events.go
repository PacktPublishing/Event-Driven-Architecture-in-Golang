package handlers

import (
	"context"

	"eda-in-golang/baskets/basketspb"
	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher[ddd.Event]) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.BasketStartedEvent,
		domain.BasketCheckedOutEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.BasketStartedEvent:
		return h.onBasketStarted(ctx, event)
	case domain.BasketCanceledEvent:
		return h.onBasketCanceled(ctx, event)
	case domain.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onBasketStarted(ctx context.Context, event ddd.Event) error {
	basket := event.Payload().(*domain.Basket)
	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel,
		ddd.NewEvent(basketspb.BasketStartedEvent, &basketspb.BasketStarted{
			Id:         basket.ID(),
			CustomerId: basket.CustomerID,
		}),
	)
}

func (h domainHandlers[T]) onBasketCanceled(ctx context.Context, event ddd.Event) error {
	basket := event.Payload().(*domain.Basket)
	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel,
		ddd.NewEvent(basketspb.BasketCanceledEvent, &basketspb.BasketCanceled{
			Id: basket.ID(),
		}),
	)
}

func (h domainHandlers[T]) onBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	basket := event.Payload().(*domain.Basket)
	items := make([]*basketspb.BasketCheckedOut_Item, 0, len(basket.Items))
	for _, item := range basket.Items {
		items = append(items, &basketspb.BasketCheckedOut_Item{
			StoreId:     item.StoreID,
			ProductId:   item.ProductID,
			StoreName:   item.StoreName,
			ProductName: item.ProductName,
			Price:       item.ProductPrice,
			Quantity:    int32(item.Quantity),
		})
	}
	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel,
		ddd.NewEvent(basketspb.BasketCheckedOutEvent, &basketspb.BasketCheckedOut{
			Id:         basket.ID(),
			CustomerId: basket.CustomerID,
			PaymentId:  basket.PaymentID,
			Items:      items,
		}),
	)
}
