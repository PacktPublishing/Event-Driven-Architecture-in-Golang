package ordering

import (
	"context"

	"eda-in-golang/internal/monolith"
	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/grpc"
	"eda-in-golang/ordering/internal/logging"
	"eda-in-golang/ordering/internal/postgres"
	"eda-in-golang/ordering/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup Driven adapters
	orders := postgres.NewOrderRepository("ordering.orders", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	customers := grpc.NewCustomerRepository(conn)
	payments := grpc.NewPaymentRepository(conn)
	invoices := grpc.NewInvoiceRepository(conn)
	shopping := grpc.NewShoppingListRepository(conn)
	notifications := grpc.NewNotificationRepository(conn)

	// setup application
	var app application.App
	app = application.New(orders, customers, payments, invoices, shopping, notifications)
	app = logging.NewApplication(app, mono.Logger())

	// setup Driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	return nil
}
