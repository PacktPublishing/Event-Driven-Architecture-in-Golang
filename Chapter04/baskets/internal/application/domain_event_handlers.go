package application

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type DomainEventHandlers interface {
	OnBasketStarted(ctx context.Context, event ddd.Event) error
	OnBasketItemAdded(ctx context.Context, event ddd.Event) error
	OnBasketItemRemoved(ctx context.Context, event ddd.Event) error
	OnBasketCanceled(ctx context.Context, event ddd.Event) error
	OnBasketCheckedOut(ctx context.Context, event ddd.Event) error
}

type ignoreUnimplementedDomainEvents struct{}

var _ DomainEventHandlers = (*ignoreUnimplementedDomainEvents)(nil)

func (ignoreUnimplementedDomainEvents) OnBasketStarted(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnBasketItemAdded(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnBasketItemRemoved(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnBasketCanceled(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	return nil
}
