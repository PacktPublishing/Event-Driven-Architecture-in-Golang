package handlers

import (
	"context"
	"database/sql"

	"eda-in-golang/customers/customerspb"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/di"
	"eda-in-golang/internal/registry"
	"eda-in-golang/ordering/orderingpb"
	"eda-in-golang/stores/storespb"
)

func RegisterIntegrationEventHandlersTx(container di.Container) (err error) {
	evtMsgHandler := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.IncomingRawMessage) (err error) {
		ctx = container.Scoped(ctx)
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}(di.Get(ctx, "tx").(*sql.Tx))

		evtHandlers := am.RawMessageHandlerWithMiddleware(
			am.NewEventMessageHandler(
				di.Get(ctx, "registry").(registry.Registry),
				di.Get(ctx, "integrationEventHandlers").(ddd.EventHandler[ddd.Event]),
			),
			di.Get(ctx, "inboxMiddleware").(am.RawMessageHandlerMiddleware),
		)

		return evtHandlers.HandleMessage(ctx, msg)
	})

	subscriber := container.Get("stream").(am.RawMessageStream)

	if err = subscriber.Subscribe(customerspb.CustomerAggregateChannel, evtMsgHandler, am.MessageFilter{
		customerspb.CustomerRegisteredEvent,
	}, am.GroupName("search-customers")); err != nil {
		return
	}

	if err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, evtMsgHandler, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("notification-orders")); err != nil {
		return
	}

	if err = subscriber.Subscribe(storespb.ProductAggregateChannel, evtMsgHandler, am.MessageFilter{
		storespb.ProductAddedEvent,
		storespb.ProductRebrandedEvent,
		storespb.ProductRemovedEvent,
	}, am.GroupName("search-products")); err != nil {
		return
	}

	if err = subscriber.Subscribe(storespb.StoreAggregateChannel, evtMsgHandler, am.MessageFilter{
		storespb.StoreCreatedEvent,
		storespb.StoreRebrandedEvent,
	}, am.GroupName("search-stores")); err != nil {
		return
	}

	return
}
