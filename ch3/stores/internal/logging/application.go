package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/stores/internal/application"
	"eda-in-golang/stores/internal/application/commands"
	"eda-in-golang/stores/internal/application/queries"
	"eda-in-golang/stores/internal/domain"
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

func (a Application) CreateStore(ctx context.Context, cmd commands.CreateStore) (err error) {
	a.logger.Info().Msg("--> Stores.CreateStore")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.CreateStore") }()
	return a.App.CreateStore(ctx, cmd)
}

func (a Application) EnableParticipation(ctx context.Context, cmd commands.EnableParticipation) (err error) {
	a.logger.Info().Msg("--> Stores.EnableParticipation")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.EnableParticipation") }()
	return a.App.EnableParticipation(ctx, cmd)
}

func (a Application) DisableParticipation(ctx context.Context, cmd commands.DisableParticipation) (err error) {
	a.logger.Info().Msg("--> Stores.DisableParticipation")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.DisableParticipation") }()
	return a.App.DisableParticipation(ctx, cmd)
}

func (a Application) AddProduct(ctx context.Context, cmd commands.AddProduct) (err error) {
	a.logger.Info().Msg("--> Stores.AddProduct")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.AddProduct") }()
	return a.App.AddProduct(ctx, cmd)
}

func (a Application) RemoveProduct(ctx context.Context, cmd commands.RemoveProduct) (err error) {
	a.logger.Info().Msg("--> Stores.RemoveProduct")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.RemoveProduct") }()
	return a.App.RemoveProduct(ctx, cmd)
}

func (a Application) GetStore(ctx context.Context, query queries.GetStore) (store *domain.Store, err error) {
	a.logger.Info().Msg("--> Stores.GetStore")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetStore") }()
	return a.App.GetStore(ctx, query)
}

func (a Application) GetStores(ctx context.Context, query queries.GetStores) (stores []*domain.Store, err error) {
	a.logger.Info().Msg("--> Stores.GetStores")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetStores") }()
	return a.App.GetStores(ctx, query)
}

func (a Application) GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStores,
) (store []*domain.Store,
	err error,
) {
	a.logger.Info().Msg("--> Stores.GetParticipatingStores")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetParticipatingStores") }()
	return a.App.GetParticipatingStores(ctx, query)
}

func (a Application) GetCatalog(ctx context.Context, query queries.GetCatalog) (products []*domain.Product, err error) {
	a.logger.Info().Msg("--> Stores.GetCatalog")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetCatalog") }()
	return a.App.GetCatalog(ctx, query)
}

func (a Application) GetProduct(ctx context.Context, query queries.GetProduct) (product *domain.Product, err error) {
	a.logger.Info().Msg("--> Stores.GetProduct")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetProduct") }()
	return a.App.GetProduct(ctx, query)
}
