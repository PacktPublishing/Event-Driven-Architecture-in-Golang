package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/models"
)

type Application struct {
	application.Application
	logger zerolog.Logger
}

var _ application.Application = (*Application)(nil)

func LogApplicationAccess(application application.Application, logger zerolog.Logger) Application {
	return Application{
		Application: application,
		logger:      logger,
	}
}

func (a Application) SearchOrders(ctx context.Context, search application.SearchOrders) (orders []*models.Order, err error) {
	a.logger.Info().Msg("--> Search.SearchOrders")
	defer func() { a.logger.Info().Err(err).Msg("<-- Search.SearchOrders") }()
	return a.Application.SearchOrders(ctx, search)
}

func (a Application) GetOrder(ctx context.Context, get application.GetOrder) (order *models.Order, err error) {
	a.logger.Info().Msg("--> Search.GetOrder")
	defer func() { a.logger.Info().Err(err).Msg("<-- Search.GetOrder") }()
	return a.Application.GetOrder(ctx, get)
}
