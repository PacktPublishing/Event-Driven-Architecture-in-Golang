package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/domain"
)

func RegisterInvoiceHandlers(invoiceHandlers application.DomainEventHandlers, domainSubscriber ddd.EventSubscriber) {
	domainSubscriber.Subscribe(domain.OrderReadied{}, invoiceHandlers.OnOrderReadied)
}
