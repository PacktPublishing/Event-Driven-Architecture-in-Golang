package application

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
	"eda-in-golang/stores/storespb"
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

func (h IntegrationEventHandlers[T]) onStoreCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreCreated)
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreCreatedEvent, &storespb.StoreCreated{
			Id:       event.ID(),
			Name:     payload.Name,
			Location: payload.Location,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onStoreParticipationEnabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreParticipatingToggledEvent, &storespb.StoreParticipationToggled{
			Id:            event.ID(),
			Participating: true,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onStoreParticipationDisabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreParticipatingToggledEvent, &storespb.StoreParticipationToggled{
			Id:            event.ID(),
			Participating: false,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreRebranded)
	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel,
		ddd.NewEvent(storespb.StoreRebrandedEvent, &storespb.StoreRebranded{
			Id:   event.ID(),
			Name: payload.Name,
		}),
	)
}
