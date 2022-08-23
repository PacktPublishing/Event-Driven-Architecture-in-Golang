package application

import (
	"context"

	"eda-in-golang/depot/internal/application/commands"
	"eda-in-golang/depot/internal/application/queries"
	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/ddd"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateShoppingList(ctx context.Context, cmd commands.CreateShoppingList) error
		CancelShoppingList(ctx context.Context, cmd commands.CancelShoppingList) error
		AssignShoppingList(ctx context.Context, cmd commands.AssignShoppingList) error
		CompleteShoppingList(ctx context.Context, cmd commands.CompleteShoppingList) error
	}
	Queries interface {
		GetShoppingList(ctx context.Context, query queries.GetShoppingList) (*domain.ShoppingList, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateShoppingListHandler
		commands.CancelShoppingListHandler
		commands.AssignShoppingListHandler
		commands.CompleteShoppingListHandler
	}
	appQueries struct {
		queries.GetShoppingListHandler
	}
)

var _ App = (*Application)(nil)

func New(shoppingLists domain.ShoppingListRepository, stores domain.StoreRepository, products domain.ProductRepository,
	domainPublisher ddd.EventPublisher[ddd.AggregateEvent],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateShoppingListHandler:   commands.NewCreateShoppingListHandler(shoppingLists, stores, products, domainPublisher),
			CancelShoppingListHandler:   commands.NewCancelShoppingListHandler(shoppingLists, domainPublisher),
			AssignShoppingListHandler:   commands.NewAssignShoppingListHandler(shoppingLists, domainPublisher),
			CompleteShoppingListHandler: commands.NewCompleteShoppingListHandler(shoppingLists, domainPublisher),
		},
		appQueries: appQueries{
			GetShoppingListHandler: queries.NewGetShoppingListHandler(shoppingLists),
		},
	}
}
