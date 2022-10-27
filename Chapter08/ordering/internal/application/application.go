package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/application/commands"
	"eda-in-golang/ordering/internal/application/queries"
	"eda-in-golang/ordering/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateOrder(ctx context.Context, cmd commands.CreateOrder) error
		RejectOrder(ctx context.Context, cmd commands.RejectOrder) error
		ApproveOrder(ctx context.Context, cmd commands.ApproveOrder) error
		CancelOrder(ctx context.Context, cmd commands.CancelOrder) error
		ReadyOrder(ctx context.Context, cmd commands.ReadyOrder) error
		CompleteOrder(ctx context.Context, cmd commands.CompleteOrder) error
	}
	Queries interface {
		GetOrder(ctx context.Context, query queries.GetOrder) (*domain.Order, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateOrderHandler
		commands.RejectOrderHandler
		commands.ApproveOrderHandler
		commands.CancelOrderHandler
		commands.ReadyOrderHandler
		commands.CompleteOrderHandler
	}
	appQueries struct {
		queries.GetOrderHandler
	}
)

var _ App = (*Application)(nil)

func New(orders domain.OrderRepository, shopping domain.ShoppingRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			CreateOrderHandler:   commands.NewCreateOrderHandler(orders, publisher),
			RejectOrderHandler:   commands.NewRejectOrderHandler(orders, publisher),
			ApproveOrderHandler:  commands.NewApproveOrderHandler(orders, publisher),
			CancelOrderHandler:   commands.NewCancelOrderHandler(orders, shopping, publisher),
			ReadyOrderHandler:    commands.NewReadyOrderHandler(orders, publisher),
			CompleteOrderHandler: commands.NewCompleteOrderHandler(orders, publisher),
		},
		appQueries: appQueries{
			GetOrderHandler: queries.NewGetOrderHandler(orders),
		},
	}
}
