package application

import (
	"context"

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
		commands.CancelOrderHandler
		commands.ReadyOrderHandler
		commands.CompleteOrderHandler
	}
	appQueries struct {
		queries.GetOrderHandler
	}
)

var _ App = (*Application)(nil)

func New(orders domain.OrderRepository, customers domain.CustomerRepository, payments domain.PaymentRepository,
	shopping domain.ShoppingRepository,
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateOrderHandler:   commands.NewCreateOrderHandler(orders, customers, payments, shopping),
			CancelOrderHandler:   commands.NewCancelOrderHandler(orders, shopping),
			ReadyOrderHandler:    commands.NewReadyOrderHandler(orders),
			CompleteOrderHandler: commands.NewCompleteOrderHandler(orders),
		},
		appQueries: appQueries{
			GetOrderHandler: queries.NewGetOrderHandler(orders),
		},
	}
}
