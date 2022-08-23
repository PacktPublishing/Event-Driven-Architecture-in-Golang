package application

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type DomainEventHandlers interface {
	OnCustomerRegistered(ctx context.Context, event ddd.Event) error
	OnCustomerAuthorized(ctx context.Context, event ddd.Event) error
	OnCustomerEnabled(ctx context.Context, event ddd.Event) error
	OnCustomerDisabled(ctx context.Context, event ddd.Event) error
}

type ignoreUnimplementedDomainEvents struct{}

var _ DomainEventHandlers = (*ignoreUnimplementedDomainEvents)(nil)

func (ignoreUnimplementedDomainEvents) OnCustomerRegistered(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnCustomerAuthorized(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnCustomerEnabled(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvents) OnCustomerDisabled(ctx context.Context, event ddd.Event) error {
	return nil
}
