package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/storespb"
)

type ProductHandlers[T ddd.Event] struct {
	cache ProductCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*ProductHandlers[ddd.Event])(nil)

func NewProductHandlers(cache ProductCacheRepository) ProductHandlers[ddd.Event] {
	return ProductHandlers[ddd.Event]{
		cache: cache,
	}
}

func (h ProductHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case storespb.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case storespb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}

	return nil
}

func (h ProductHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductAdded)
	return h.cache.Add(ctx, payload.GetId(), payload.GetStoreId(), payload.GetName())
}

func (h ProductHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRebranded)
	return h.cache.Rebrand(ctx, payload.GetId(), payload.GetName())
}

func (h ProductHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRemoved)
	return h.cache.Remove(ctx, payload.GetId())
}
