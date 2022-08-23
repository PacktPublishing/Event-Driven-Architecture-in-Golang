package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/application/commands"
	"eda-in-golang/stores/internal/application/queries"
	"eda-in-golang/stores/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateStore(ctx context.Context, cmd commands.CreateStore) error
		EnableParticipation(ctx context.Context, cmd commands.EnableParticipation) error
		DisableParticipation(ctx context.Context, cmd commands.DisableParticipation) error
		AddProduct(ctx context.Context, cmd commands.AddProduct) error
		RemoveProduct(ctx context.Context, cmd commands.RemoveProduct) error
	}
	Queries interface {
		GetStore(ctx context.Context, query queries.GetStore) (*domain.Store, error)
		GetStores(ctx context.Context, query queries.GetStores) ([]*domain.Store, error)
		GetParticipatingStores(ctx context.Context, query queries.GetParticipatingStores) ([]*domain.Store, error)
		GetCatalog(ctx context.Context, query queries.GetCatalog) ([]*domain.Product, error)
		GetProduct(ctx context.Context, query queries.GetProduct) (*domain.Product, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateStoreHandler
		commands.EnableParticipationHandler
		commands.DisableParticipationHandler
		commands.AddProductHandler
		commands.RemoveProductHandler
	}
	appQueries struct {
		queries.GetStoreHandler
		queries.GetStoresHandler
		queries.GetParticipatingStoresHandler
		queries.GetCatalogHandler
		queries.GetProductHandler
	}
)

var _ App = (*Application)(nil)

func New(stores domain.StoreRepository, participatingStores domain.ParticipatingStoreRepository,
	products domain.ProductRepository, domainPublisher ddd.EventPublisher,
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateStoreHandler:          commands.NewCreateStoreHandler(stores, domainPublisher),
			EnableParticipationHandler:  commands.NewEnableParticipationHandler(stores, domainPublisher),
			DisableParticipationHandler: commands.NewDisableParticipationHandler(stores, domainPublisher),
			AddProductHandler:           commands.NewAddProductHandler(stores, products, domainPublisher),
			RemoveProductHandler:        commands.NewRemoveProductHandler(products, domainPublisher),
		},
		appQueries: appQueries{
			GetStoreHandler:               queries.NewGetStoreHandler(stores),
			GetStoresHandler:              queries.NewGetStoresHandler(stores),
			GetParticipatingStoresHandler: queries.NewGetParticipatingStoresHandler(participatingStores),
			GetCatalogHandler:             queries.NewGetCatalogHandler(products),
			GetProductHandler:             queries.NewGetProductHandler(products),
		},
	}
}
