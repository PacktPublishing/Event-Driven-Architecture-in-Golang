package application

import (
	"context"

	"eda-in-golang/customers/customerspb"
	"eda-in-golang/customers/internal/domain"
	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type IntegrationEventHandlers[T ddd.AggregateEvent] struct {
	publisher am.MessagePublisher[ddd.Event]
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*IntegrationEventHandlers[ddd.AggregateEvent])(nil)

func NewIntegrationEventHandlers(publisher am.MessagePublisher[ddd.Event]) *IntegrationEventHandlers[ddd.AggregateEvent] {
	return &IntegrationEventHandlers[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func (h IntegrationEventHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	case domain.CustomerSmsChangedEvent:
		return h.onCustomerSmsChanged(ctx, event)
	case domain.CustomerEnabledEvent:
		return h.onCustomerEnabled(ctx, event)
	case domain.CustomerDisabledEvent:
		return h.onCustomerDisabled(ctx, event)
	}
	return nil
}

func (h IntegrationEventHandlers[T]) onCustomerRegistered(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.CustomerRegistered)
	return h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel,
		ddd.NewEvent(customerspb.CustomerRegisteredEvent, &customerspb.CustomerRegistered{
			Id:        payload.Customer.ID(),
			Name:      payload.Customer.Name,
			SmsNumber: payload.Customer.SmsNumber,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onCustomerSmsChanged(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.CustomerRegistered)
	return h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel,
		ddd.NewEvent(customerspb.CustomerSmsChangedEvent, &customerspb.CustomerSmsChanged{
			Id:        payload.Customer.ID(),
			SmsNumber: payload.Customer.SmsNumber,
		}),
	)
}

func (h IntegrationEventHandlers[T]) onCustomerEnabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel,
		ddd.NewEvent(customerspb.CustomerEnabledEvent, &customerspb.CustomerEnabled{
			Id: event.AggregateID(),
		}),
	)
}

func (h IntegrationEventHandlers[T]) onCustomerDisabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.publisher.Publish(ctx, customerspb.CustomerAggregateChannel,
		ddd.NewEvent(customerspb.CustomerDisabledEvent, &customerspb.CustomerDisabled{
			Id: event.AggregateID(),
		}),
	)
}
