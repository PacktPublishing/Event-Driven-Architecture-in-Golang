package application

import (
	"context"

	"eda-in-golang/depot/internal/application/commands"
	"eda-in-golang/depot/internal/application/queries"
	"eda-in-golang/depot/internal/domain"
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

func New(shoppingLists domain.ShoppingListRepository, stores domain.StoreRepository, products domain.ProductRepository, orders domain.OrderRepository) *Application {
	return &Application{
		appCommands: appCommands{
			CreateShoppingListHandler:   commands.NewCreateShoppingListHandler(shoppingLists, stores, products),
			CancelShoppingListHandler:   commands.NewCancelShoppingListHandler(shoppingLists),
			AssignShoppingListHandler:   commands.NewAssignShoppingListHandler(shoppingLists),
			CompleteShoppingListHandler: commands.NewCompleteShoppingListHandler(shoppingLists, orders),
		},
		appQueries: appQueries{
			GetShoppingListHandler: queries.NewGetShoppingListHandler(shoppingLists),
		},
	}
}
