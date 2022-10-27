package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/payments/internal/application"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(application application.App, logger zerolog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) AuthorizePayment(ctx context.Context, authorize application.AuthorizePayment) (err error) {
	a.logger.Info().Msg("--> Payments.AuthorizePayment")
	defer func() { a.logger.Info().Err(err).Msg("<-- Payments.AuthorizePayment") }()
	return a.App.AuthorizePayment(ctx, authorize)
}

func (a Application) ConfirmPayment(ctx context.Context, confirm application.ConfirmPayment) (err error) {
	a.logger.Info().Msg("--> Payments.ConfirmPayment")
	defer func() { a.logger.Info().Err(err).Msg("<-- Payments.ConfirmPayment") }()
	return a.App.ConfirmPayment(ctx, confirm)
}

func (a Application) CreateInvoice(ctx context.Context, create application.CreateInvoice) (err error) {
	a.logger.Info().Msg("--> Payments.CreateInvoice")
	defer func() { a.logger.Info().Err(err).Msg("<-- Payments.CreateInvoice") }()
	return a.App.CreateInvoice(ctx, create)
}

func (a Application) AdjustInvoice(ctx context.Context, adjust application.AdjustInvoice) (err error) {
	a.logger.Info().Msg("--> Payments.AdjustInvoice")
	defer func() { a.logger.Info().Err(err).Msg("<-- Payments.AdjustInvoice") }()
	return a.App.AdjustInvoice(ctx, adjust)
}

func (a Application) PayInvoice(ctx context.Context, pay application.PayInvoice) (err error) {
	a.logger.Info().Msg("--> Payments.PayInvoice")
	defer func() { a.logger.Info().Err(err).Msg("<-- Payments.PayInvoice") }()
	return a.App.PayInvoice(ctx, pay)
}

func (a Application) CancelInvoice(ctx context.Context, cancel application.CancelInvoice) (err error) {
	a.logger.Info().Msg("--> Payments.CancelInvoice")
	defer func() { a.logger.Info().Err(err).Msg("<-- Payments.CancelInvoice") }()
	return a.App.CancelInvoice(ctx, cancel)
}
