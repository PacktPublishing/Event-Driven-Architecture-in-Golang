package notifications

import (
	"context"

	"eda-in-golang/customers/customerspb"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/jetstream"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/internal/registry"
	"eda-in-golang/notifications/internal/application"
	"eda-in-golang/notifications/internal/grpc"
	"eda-in-golang/notifications/internal/handlers"
	"eda-in-golang/notifications/internal/logging"
	"eda-in-golang/notifications/internal/postgres"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = customerspb.Registrations(reg); err != nil {
		return err
	}
	eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	customers := postgres.NewCustomerCacheRepository("notifications.customers_cache", mono.DB(), grpc.NewCustomerRepository(conn))

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers),
		mono.Logger(),
	)
	customerHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewCustomerHandlers(customers),
		"Customer", mono.Logger(),
	)
	orderHandlers := logging.LogEventHandlerAccess[ddd.Event](
		application.NewOrderHandlers(app),
		"Order", mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err = handlers.RegisterCustomerHandlers(customerHandlers, eventStream); err != nil {
		return err
	}
	if err = handlers.RegisterOrderHandlers(orderHandlers, eventStream); err != nil {
		return err
	}

	return nil
}
