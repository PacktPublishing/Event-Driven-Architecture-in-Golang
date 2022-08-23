package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/depot/internal/application"
	"eda-in-golang/internal/ddd"
)

type DomainEventHandlers struct {
	application.DomainEventHandlers
	logger zerolog.Logger
}

var _ application.DomainEventHandlers = (*DomainEventHandlers)(nil)

func LogDomainEventHandlerAccess(handlers application.DomainEventHandlers, logger zerolog.Logger) DomainEventHandlers {
	return DomainEventHandlers{
		DomainEventHandlers: handlers,
		logger:              logger,
	}
}

func (h DomainEventHandlers) OnShoppingListCreated(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Depot.OnShoppingListCreated")
	defer func() { h.logger.Info().Err(err).Msg("<-- Depot.OnShoppingListCreated") }()
	return h.DomainEventHandlers.OnShoppingListCreated(ctx, event)
}

func (h DomainEventHandlers) OnShoppingListCanceled(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Depot.OnShoppingListCanceled")
	defer func() { h.logger.Info().Err(err).Msg("<-- Depot.OnShoppingListCanceled") }()
	return h.DomainEventHandlers.OnShoppingListCanceled(ctx, event)
}

func (h DomainEventHandlers) OnShoppingListAssigned(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Depot.OnShoppingListAssigned")
	defer func() { h.logger.Info().Err(err).Msg("<-- Depot.OnShoppingListAssigned") }()
	return h.DomainEventHandlers.OnShoppingListAssigned(ctx, event)
}

func (h DomainEventHandlers) OnShoppingListCompleted(ctx context.Context, event ddd.Event) (err error) {
	h.logger.Info().Msg("--> Depot.OnShoppingListCompleted")
	defer func() { h.logger.Info().Err(err).Msg("<-- Depot.OnShoppingListCompleted") }()
	return h.DomainEventHandlers.OnShoppingListCompleted(ctx, event)
}
