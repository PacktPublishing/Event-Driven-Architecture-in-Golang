package application

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/payments/internal/models"
	"eda-in-golang/payments/paymentspb"
)

type IntegrationEventHandlers[T ddd.Event] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.Event] = (*IntegrationEventHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(publisher am.MessagePublisher[ddd.Event]) *IntegrationEventHandlers[ddd.Event] {
	return &IntegrationEventHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func (h IntegrationEventHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case models.InvoicePaidEvent:
		return h.onInvoicePaid(ctx, event)
	}
	return nil
}

func (h IntegrationEventHandlers[T]) onInvoicePaid(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*models.InvoicePaid)
	return h.publisher.Publish(ctx, paymentspb.InvoiceAggregateChannel,
		ddd.NewEvent(paymentspb.InvoicePaidEvent, &paymentspb.InvoicePaid{
			Id:      payload.ID,
			OrderId: payload.OrderID,
		}),
	)
}
