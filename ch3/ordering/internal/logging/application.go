package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/application/commands"
	"eda-in-golang/ordering/internal/application/queries"
	"eda-in-golang/ordering/internal/domain"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

var _ application.App = (*Application)(nil)

func NewApplication(application application.App, logger zerolog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) CreateOrder(ctx context.Context, cmd commands.CreateOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.CreateOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.CreateOrder") }()
	return a.App.CreateOrder(ctx, cmd)
}

func (a Application) CancelOrder(ctx context.Context, cmd commands.CancelOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.CancelOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.CancelOrder") }()
	return a.App.CancelOrder(ctx, cmd)
}

func (a Application) ReadyOrder(ctx context.Context, cmd commands.ReadyOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.ReadyOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.ReadyOrder") }()
	return a.App.ReadyOrder(ctx, cmd)
}

func (a Application) CompleteOrder(ctx context.Context, cmd commands.CompleteOrder) (err error) {
	a.logger.Info().Msg("--> Ordering.CompleteOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.CompleteOrder") }()
	return a.App.CompleteOrder(ctx, cmd)
}

func (a Application) GetOrder(ctx context.Context, query queries.GetOrder) (order *domain.Order, err error) {
	a.logger.Info().Msg("--> Ordering.GetOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Ordering.GetOrder") }()
	return a.App.GetOrder(ctx, query)
}
