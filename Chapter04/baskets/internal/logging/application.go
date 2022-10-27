package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/baskets/internal/application"
	"eda-in-golang/baskets/internal/domain"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(application application.App, logger zerolog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) StartBasket(ctx context.Context, start application.StartBasket) (err error) {
	a.logger.Info().Msg("--> Baskets.StartBasket")
	defer func() { a.logger.Info().Err(err).Msg("<-- Baskets.StartBasket") }()
	return a.App.StartBasket(ctx, start)
}

func (a Application) CancelBasket(ctx context.Context, cancel application.CancelBasket) (err error) {
	a.logger.Info().Msg("--> Baskets.CancelBasket")
	defer func() { a.logger.Info().Err(err).Msg("<-- Baskets.CancelBasket") }()
	return a.App.CancelBasket(ctx, cancel)
}

func (a Application) CheckoutBasket(ctx context.Context, checkout application.CheckoutBasket) (err error) {
	a.logger.Info().Msg("--> Baskets.CheckoutBasket")
	defer func() { a.logger.Info().Err(err).Msg("<-- Baskets.CheckoutBasket") }()
	return a.App.CheckoutBasket(ctx, checkout)
}

func (a Application) AddItem(ctx context.Context, add application.AddItem) (err error) {
	a.logger.Info().Msg("--> Baskets.AddItem")
	defer func() { a.logger.Info().Err(err).Msg("<-- Baskets.AddItem") }()
	return a.App.AddItem(ctx, add)
}

func (a Application) RemoveItem(ctx context.Context, remove application.RemoveItem) (err error) {
	a.logger.Info().Msg("--> Baskets.RemoveItem")
	defer func() { a.logger.Info().Err(err).Msg("<-- Baskets.RemoveItem") }()
	return a.App.RemoveItem(ctx, remove)
}

func (a Application) GetBasket(ctx context.Context, get application.GetBasket) (basket *domain.Basket, err error) {
	a.logger.Info().Msg("--> Baskets.GetBasket")
	defer func() { a.logger.Info().Err(err).Msg("<-- Baskets.GetBasket") }()
	return a.App.GetBasket(ctx, get)
}
