package handlers

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/errorsotel"
	"eda-in-golang/stores/internal/constants"
	"eda-in-golang/stores/internal/domain"
)

type mallHandlers[T ddd.Event] struct {
	mall domain.MallRepository
}

var _ ddd.EventHandler[ddd.Event] = (*mallHandlers[ddd.Event])(nil)

func NewMallHandlers(mall domain.MallRepository) ddd.EventHandler[ddd.Event] {
	return mallHandlers[ddd.Event]{
		mall: mall,
	}
}

func RegisterMallHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.StoreCreatedEvent,
		domain.StoreParticipationEnabledEvent,
		domain.StoreParticipationDisabledEvent,
		domain.StoreRebrandedEvent,
	)
}

func RegisterMallHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		mallHandlers := di.Get(ctx, constants.MallHandlersKey).(ddd.EventHandler[ddd.Event])

		return mallHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMallHandlers(subscriber, handlers)
}

func (h mallHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling mall event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled mall event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling mall event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case domain.StoreParticipationEnabledEvent:
		return h.onStoreParticipationEnabled(ctx, event)
	case domain.StoreParticipationDisabledEvent:
		return h.onStoreParticipationDisabled(ctx, event)
	case domain.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}
	return nil
}

func (h mallHandlers[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Store)
	return h.mall.AddStore(ctx, payload.ID(), payload.Name, payload.Location)
}

func (h mallHandlers[T]) onStoreParticipationEnabled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Store)
	return h.mall.SetStoreParticipation(ctx, payload.ID(), true)
}

func (h mallHandlers[T]) onStoreParticipationDisabled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Store)
	return h.mall.SetStoreParticipation(ctx, payload.ID(), false)
}

func (h mallHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Store)
	return h.mall.RenameStore(ctx, payload.ID(), payload.Name)
}
