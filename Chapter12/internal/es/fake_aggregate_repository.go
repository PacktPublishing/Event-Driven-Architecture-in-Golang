package es

import (
	"context"
	"fmt"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

type FakeAggregateRepository[T EventSourcedAggregate] struct {
	aggregateName string
	registry      registry.Registry
	aggregates    map[string]T
}

var _ AggregateRepository[EventSourcedAggregate] = (*FakeAggregateRepository[EventSourcedAggregate])(nil)

func NewFakeAggregateRepository[T EventSourcedAggregate](aggregateName string, registry registry.Registry) *FakeAggregateRepository[T] {
	return &FakeAggregateRepository[T]{
		aggregateName: aggregateName,
		registry:      registry,
		aggregates:    make(map[string]T),
	}
}

func (r *FakeAggregateRepository[T]) Load(ctx context.Context, aggregateID string) (agg T, err error) {
	var exists bool

	if agg, exists = r.aggregates[aggregateID]; exists {
		return agg, nil
	}

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

	return agg, nil
}

func (r *FakeAggregateRepository[T]) Save(ctx context.Context, aggregate T) error {
	r.aggregates[aggregate.ID()] = aggregate

	return nil
}

func (r *FakeAggregateRepository[T]) Reset(aggregates ...T) {
	r.aggregates = make(map[string]T)

	for _, aggregate := range aggregates {
		r.aggregates[aggregate.ID()] = aggregate
	}
}
