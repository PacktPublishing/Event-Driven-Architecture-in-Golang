package customers

import (
	"context"

	"eda-in-golang/customers/internal/application"
	"eda-in-golang/customers/internal/grpc"
	"eda-in-golang/customers/internal/logging"
	"eda-in-golang/customers/internal/postgres"
	"eda-in-golang/customers/internal/rest"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/monolith"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	customers := postgres.NewCustomerRepository("customers.customers", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers, domainDispatcher),
		mono.Logger(),
	)

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
