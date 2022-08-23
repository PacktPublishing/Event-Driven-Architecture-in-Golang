package es

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type EventPublisher struct {
	AggregateStore
	publisher ddd.EventPublisher[ddd.AggregateEvent]
}

var _ AggregateStore = (*EventPublisher)(nil)

func NewEventPublisher(publisher ddd.EventPublisher[ddd.AggregateEvent]) AggregateStoreMiddleware {
	eventPublisher := EventPublisher{
		publisher: publisher,
	}

	return func(store AggregateStore) AggregateStore {
		eventPublisher.AggregateStore = store
		return eventPublisher
	}
}

func (p EventPublisher) Save(ctx context.Context, aggregate EventSourcedAggregate) error {
	if err := p.AggregateStore.Save(ctx, aggregate); err != nil {
		return err
	}
	return p.publisher.Publish(ctx, aggregate.Events()...)
}
