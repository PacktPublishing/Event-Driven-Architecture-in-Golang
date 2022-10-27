package payments

import (
	"context"

	"eda-in-golang/internal/monolith"
	"eda-in-golang/payments/internal/application"
	"eda-in-golang/payments/internal/grpc"
	"eda-in-golang/payments/internal/logging"
	"eda-in-golang/payments/internal/postgres"
	"eda-in-golang/payments/internal/rest"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup Driven adapters
	invoices := postgres.NewInvoiceRepository("payments.invoices", mono.DB())
	payments := postgres.NewPaymentRepository("payments.payments", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	orders := grpc.NewOrderRepository(conn)

	// setup application
	var app application.App
	app = application.New(invoices, payments, orders)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
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
