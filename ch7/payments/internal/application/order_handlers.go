package application

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/orderingpb"
)

type OrderHandlers[T ddd.Event] struct {
	app App
}

var _ ddd.EventHandler[ddd.Event] = (*OrderHandlers[ddd.Event])(nil)

func NewOrderHandlers(app App) OrderHandlers[ddd.Event] {
	return OrderHandlers[ddd.Event]{
		app: app,
	}
}

func (h OrderHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case orderingpb.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	}
	return nil
}

func (h OrderHandlers[T]) onOrderReadied(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderReadied)
	return h.app.CreateInvoice(ctx, CreateInvoice{
		ID:        payload.GetId(),
		OrderID:   payload.GetId(),
		PaymentID: payload.GetPaymentId(),
		Amount:    payload.GetTotal(),
	})
}

func (h OrderHandlers[T]) onOrderCanceled(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)
	return h.app.CancelInvoice(ctx, CancelInvoice{
		ID: payload.GetId(),
	})
}
