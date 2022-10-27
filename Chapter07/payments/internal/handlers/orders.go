package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/orderingpb"
)

func RegisterOrderHandlers(orderHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, eventMsg am.EventMessage) error {
		return orderHandlers.HandleEvent(ctx, eventMsg)
	})

	return stream.Subscribe(orderingpb.OrderAggregateChannel, evtMsgHandler, am.MessageFilter{
		orderingpb.OrderReadiedEvent,
	}, am.GroupName("payment-orders"))
}
