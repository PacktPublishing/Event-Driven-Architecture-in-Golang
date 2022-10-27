package handlers

import (
	"context"

	"eda-in-golang/cosec/internal/models"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/sec"
	"eda-in-golang/ordering/orderingpb"
)

type integrationHandlers[T ddd.Event] struct {
	orchestrator sec.Orchestrator[*models.CreateOrderData]
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(saga sec.Orchestrator[*models.CreateOrderData]) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		orchestrator: saga,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	_, err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, evtMsgHandler, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
	}, am.GroupName("cosec-ordering"))
	return
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*orderingpb.OrderCreated)

	var total float64
	items := make([]models.Item, len(payload.GetItems()))
	for i, item := range payload.GetItems() {
		items[i] = models.Item{
			ProductID: item.GetProductId(),
			StoreID:   item.GetStoreId(),
			Price:     item.GetPrice(),
			Quantity:  int(item.GetQuantity()),
		}
		total += float64(item.GetQuantity()) * item.GetPrice()
	}

	data := &models.CreateOrderData{
		OrderID:    payload.GetId(),
		CustomerID: payload.GetCustomerId(),
		PaymentID:  payload.GetPaymentId(),
		Items:      items,
		Total:      total,
	}

	// Start the CreateOrderSaga
	return h.orchestrator.Start(ctx, event.ID(), data)
}
