package es

import (
	"context"
	"fmt"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

type (
	AggregateRepository[T EventSourcedAggregate] interface {
		Load(ctx context.Context, aggregateID string) (agg T, err error)
		Save(ctx context.Context, aggregate T) error
	}

	aggregateRepository[T EventSourcedAggregate] struct {
		aggregateName string
		registry      registry.Registry
		store         AggregateStore
	}
)

var _ AggregateRepository[EventSourcedAggregate] = (*aggregateRepository[EventSourcedAggregate])(nil)

func NewAggregateRepository[T EventSourcedAggregate](aggregateName string, registry registry.Registry, store AggregateStore) aggregateRepository[T] {
	return aggregateRepository[T]{
		aggregateName: aggregateName,
		registry:      registry,
		store:         store,
	}
}

func (r aggregateRepository[T]) Load(ctx context.Context, aggregateID string) (agg T, err error) {
	var v any
	v, err = r.registry.Build(
		r.aggregateName,
		ddd.SetID(aggregateID),
		ddd.SetName(r.aggregateName),
	)
	if err != nil {
		return agg, err
	}

	var ok bool
	if agg, ok = v.(T); !ok {
		return agg, fmt.Errorf("%T is not the expected type %T", v, agg)
	}

	if err = r.store.Load(ctx, agg); err != nil {
		return agg, err
	}

	return agg, nil
}

func (r aggregateRepository[T]) Save(ctx context.Context, aggregate T) error {
	if aggregate.Version() == aggregate.PendingVersion() {
		return nil
	}

	for _, event := range aggregate.Events() {
		if err := aggregate.ApplyEvent(event); err != nil {
			return err
		}
	}

	err := r.store.Save(ctx, aggregate)
	if err != nil {
		return err
	}

	aggregate.CommitEvents()

	return nil
}
