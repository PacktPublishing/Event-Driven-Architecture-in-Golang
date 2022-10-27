package application

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/storespb"
)

type StoreHandlers[T ddd.Event] struct {
	logger zerolog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*StoreHandlers[ddd.Event])(nil)

func NewStoreHandlers(logger zerolog.Logger) StoreHandlers[ddd.Event] {
	return StoreHandlers[ddd.Event]{
		logger: logger,
	}
}

func (h StoreHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storespb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storespb.StoreParticipatingToggledEvent:
		return h.onStoreParticipationToggled(ctx, event)
	case storespb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}

	return nil
}

func (h StoreHandlers[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreCreated)
	h.logger.Debug().Msgf(`ID: %s, Name: "%s", Location: "%s"`, payload.GetId(), payload.GetName(), payload.GetLocation())
	return nil
}

func (h StoreHandlers[T]) onStoreParticipationToggled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreParticipationToggled)
	h.logger.Debug().Msgf(`ID: %s, Participating: %b`, payload.GetId(), payload.Participating)
	return nil
}

func (h StoreHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storespb.StoreRebranded)
	h.logger.Debug().Msgf(`ID: %s, Name: "%s"`, payload.GetId(), payload.GetName())
	return nil
}
