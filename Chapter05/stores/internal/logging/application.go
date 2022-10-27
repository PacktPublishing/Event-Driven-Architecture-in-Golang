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

func (a Application) RebrandStore(ctx context.Context, cmd commands.RebrandStore) (err error) {
	a.logger.Info().Msg("--> Stores.RebrandStore")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.RebrandStore") }()
	return a.App.RebrandStore(ctx, cmd)
}

func (a Application) AddProduct(ctx context.Context, cmd commands.AddProduct) (err error) {
	a.logger.Info().Msg("--> Stores.AddProduct")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.AddProduct") }()
	return a.App.AddProduct(ctx, cmd)
}

func (a Application) RebrandProduct(ctx context.Context, cmd commands.RebrandProduct) (err error) {
	a.logger.Info().Msg("--> Products.RebrandProduct")
	defer func() { a.logger.Info().Err(err).Msg("<-- Products.RebrandProduct") }()
	return a.App.RebrandProduct(ctx, cmd)
}

func (a Application) IncreaseProductPrice(ctx context.Context, cmd commands.IncreaseProductPrice) (err error) {
	a.logger.Info().Msg("--> Products.IncreaseProductPrice")
	defer func() { a.logger.Info().Err(err).Msg("<-- Products.IncreaseProductPrice") }()
	return a.App.IncreaseProductPrice(ctx, cmd)
}

func (a Application) DecreaseProductPrice(ctx context.Context, cmd commands.DecreaseProductPrice) (err error) {
	a.logger.Info().Msg("--> Products.DecreaseProductPrice")
	defer func() { a.logger.Info().Err(err).Msg("<-- Products.DecreaseProductPrice") }()
	return a.App.DecreaseProductPrice(ctx, cmd)
}

func (a Application) RemoveProduct(ctx context.Context, cmd commands.RemoveProduct) (err error) {
	a.logger.Info().Msg("--> Stores.RemoveProduct")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.RemoveProduct") }()
	return a.App.RemoveProduct(ctx, cmd)
}

func (a Application) GetStore(ctx context.Context, query queries.GetStore) (store *domain.MallStore, err error) {
	a.logger.Info().Msg("--> Stores.GetStore")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetStore") }()
	return a.App.GetStore(ctx, query)
}

func (a Application) GetStores(ctx context.Context, query queries.GetStores) (stores []*domain.MallStore, err error) {
	a.logger.Info().Msg("--> Stores.GetStores")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetStores") }()
	return a.App.GetStores(ctx, query)
}

func (a Application) GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStores) (store []*domain.MallStore, err error) {
	a.logger.Info().Msg("--> Stores.GetParticipatingStores")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetParticipatingStores") }()
	return a.App.GetParticipatingStores(ctx, query)
}

func (a Application) GetCatalog(ctx context.Context, query queries.GetCatalog) (products []*domain.CatalogProduct, err error) {
	a.logger.Info().Msg("--> Stores.GetCatalog")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetCatalog") }()
	return a.App.GetCatalog(ctx, query)
}

func (a Application) GetProduct(ctx context.Context, query queries.GetProduct) (product *domain.CatalogProduct, err error) {
	a.logger.Info().Msg("--> Stores.GetProduct")
	defer func() { a.logger.Info().Err(err).Msg("<-- Stores.GetProduct") }()
	return a.App.GetProduct(ctx, query)
}
