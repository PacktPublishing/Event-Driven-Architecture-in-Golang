package application

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/storespb"
)

type ProductHandlers[T ddd.Event] struct {
	logger zerolog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*ProductHandlers[ddd.Event])(nil)

func NewProductHandlers(logger zerolog.Logger) ProductHandlers[ddd.Event] {
	return ProductHandlers[ddd.Event]{
		logger: logger,
	}
}

func (h ProductHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case storespb.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case storespb.ProductPriceIncreasedEvent, storespb.ProductPriceDecreasedEvent:
		return h.onProductPriceChanged(ctx, event)
	case storespb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}

	return nil
}

func (h ProductHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductAdded)
	h.logger.Debug().Msgf(`ID: %s, Name: "%s", Price: "%d"`, payload.GetId(), payload.GetName(), payload.GetPrice())
	return nil
}

func (h ProductHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRebranded)
	h.logger.Debug().Msgf(`ID: %s, Name: "%s", Description: "%s"`, payload.GetId(), payload.GetName(), payload.GetDescription())
	return nil
}

func (h ProductHandlers[T]) onProductPriceChanged(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductPriceChanged)
	h.logger.Debug().Msgf(`ID: %s, Delta: "%d"`, payload.GetId(), payload.GetDelta())
	return nil
}

func (h ProductHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.ProductRemoved)
	h.logger.Debug().Msgf(`ID: %s, Price: "%d"`, payload.GetId())
	return nil
}
