package handlers

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/stores/internal/domain"
)

type catalogHandlers[T ddd.AggregateEvent] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*catalogHandlers[ddd.AggregateEvent])(nil)

func NewCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.AggregateEvent] {
	return catalogHandlers[ddd.AggregateEvent]{
		catalog: catalog,
	}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.AggregateEvent], handlers ddd.EventHandler[ddd.AggregateEvent]) {
	subscriber.Subscribe(handlers,
		domain.ProductAddedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreasedEvent,
		domain.ProductPriceDecreasedEvent,
		domain.ProductRemovedEvent,
	)
}

func RegisterCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {
		catalogHandlers := di.Get(ctx, "catalogHandlers").(ddd.EventHandler[ddd.AggregateEvent])

		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.AggregateEvent])

	RegisterCatalogHandlers(subscriber, handlers)
}

func (h catalogHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case domain.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case domain.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)
	case domain.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)
	case domain.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}
	return nil
}

func (h catalogHandlers[T]) onProductAdded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductAdded)
	return h.catalog.AddProduct(ctx, event.AggregateID(), payload.StoreID, payload.Name, payload.Description, payload.SKU, payload.Price)
}

func (h catalogHandlers[T]) onProductRebranded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductRebranded)
	return h.catalog.Rebrand(ctx, event.AggregateID(), payload.Name, payload.Description)
}

func (h catalogHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductPriceChanged)
	return h.catalog.UpdatePrice(ctx, event.AggregateID(), payload.Delta)
}

func (h catalogHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.ProductPriceChanged)
	return h.catalog.UpdatePrice(ctx, event.AggregateID(), payload.Delta)
}

func (h catalogHandlers[T]) onProductRemoved(ctx context.Context, event ddd.AggregateEvent) error {
	return h.catalog.RemoveProduct(ctx, event.AggregateID())
}
