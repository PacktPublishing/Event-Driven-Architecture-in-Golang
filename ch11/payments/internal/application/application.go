package application

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/payments/internal/models"
)

type (
	AuthorizePayment struct {
		ID         string
		CustomerID string
		Amount     float64
	}

	ConfirmPayment struct {
		ID string
	}

	CreateInvoice struct {
		ID        string
		OrderID   string
		PaymentID string
		Amount    float64
	}

	AdjustInvoice struct {
		ID     string
		Amount float64
	}

	PayInvoice struct {
		ID string
	}

	CancelInvoice struct {
		ID string
	}

	App interface {
		AuthorizePayment(ctx context.Context, authorize AuthorizePayment) error
		ConfirmPayment(ctx context.Context, confirm ConfirmPayment) error
		CreateInvoice(ctx context.Context, create CreateInvoice) error
		AdjustInvoice(ctx context.Context, adjust AdjustInvoice) error
		PayInvoice(ctx context.Context, pay PayInvoice) error
		CancelInvoice(ctx context.Context, cancel CancelInvoice) error
	}

	Application struct {
		invoices  InvoiceRepository
		payments  PaymentRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

var _ App = (*Application)(nil)

func New(invoices InvoiceRepository, payments PaymentRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		invoices:  invoices,
		payments:  payments,
		publisher: publisher,
	}
}

func (a Application) AuthorizePayment(ctx context.Context, authorize AuthorizePayment) error {
	return a.payments.Save(ctx, &models.Payment{
		ID:         authorize.ID,
		CustomerID: authorize.CustomerID,
		Amount:     authorize.Amount,
	})
}

func (a Application) ConfirmPayment(ctx context.Context, confirm ConfirmPayment) error {
	if payment, err := a.payments.Find(ctx, confirm.ID); err != nil || payment == nil {
		return errors.Wrap(errors.ErrNotFound, "payment cannot be confirmed")
	}

	return nil
}

func (a Application) CreateInvoice(ctx context.Context, create CreateInvoice) error {
	return a.invoices.Save(ctx, &models.Invoice{
		ID:      create.ID,
		OrderID: create.OrderID,
		Amount:  create.Amount,
		Status:  models.InvoiceIsPending,
	})
}

func (a Application) AdjustInvoice(ctx context.Context, adjust AdjustInvoice) error {
	invoice, err := a.invoices.Find(ctx, adjust.ID)
	if err != nil {
		return err
	}

	invoice.Amount = adjust.Amount

	return a.invoices.Update(ctx, invoice)
}

func (a Application) PayInvoice(ctx context.Context, pay PayInvoice) error {
	invoice, err := a.invoices.Find(ctx, pay.ID)
	if err != nil {
		return err
	}

	if invoice.Status != models.InvoiceIsPending {
		return errors.Wrap(errors.ErrBadRequest, "invoice cannot be paid for")
	}

	invoice.Status = models.InvoiceIsPaid

	// Before or after the invoice is saved we still risk something failing which
	// will leave the state change only partially complete
	if err = a.publisher.Publish(ctx, ddd.NewEvent(models.InvoicePaidEvent, &models.InvoicePaid{
		ID:      invoice.ID,
		OrderID: invoice.OrderID,
	})); err != nil {
		return err
	}

	return a.invoices.Update(ctx, invoice)
}

func (a Application) CancelInvoice(ctx context.Context, cancel CancelInvoice) error {
	invoice, err := a.invoices.Find(ctx, cancel.ID)
	if err != nil {
		return err
	}

	if invoice.Status != models.InvoiceIsPending {
		return errors.Wrap(errors.ErrBadRequest, "invoice cannot be paid for")
	}

	invoice.Status = models.InvoiceIsCanceled

	return a.invoices.Update(ctx, invoice)
}
