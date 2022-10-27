package application

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type DomainEventHandlers interface {
	OnOrderCreated(ctx context.Context, event ddd.Event) error
	OnOrderReadied(ctx context.Context, event ddd.Event) error
	OnOrderCanceled(ctx context.Context, event ddd.Event) error
	OnOrderCompleted(ctx context.Context, event ddd.Event) error
}

type ignoreUnimplementedDomainEvents struct{}

var _ DomainEventHandlers = (*ignoreUnimplementedDomainEvents)(nil)

func (ignoreUnimplementedDomainEvents) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnOrderCanceled(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnOrderCompleted(ctx context.Context, event ddd.Event) error {
	return nil
}
