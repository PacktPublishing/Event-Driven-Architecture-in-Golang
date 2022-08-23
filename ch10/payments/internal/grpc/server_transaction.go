package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc"

	"eda-in-golang/internal/di"
	"eda-in-golang/payments/internal/application"
	"eda-in-golang/payments/paymentspb"
)

type serverTx struct {
	c di.Container
	paymentspb.UnimplementedPaymentsServiceServer
}

var _ paymentspb.PaymentsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	paymentspb.RegisterPaymentsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) AuthorizePayment(ctx context.Context, request *paymentspb.AuthorizePaymentRequest) (resp *paymentspb.AuthorizePaymentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.AuthorizePayment(ctx, request)
}

func (s serverTx) ConfirmPayment(ctx context.Context, request *paymentspb.ConfirmPaymentRequest) (resp *paymentspb.ConfirmPaymentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.ConfirmPayment(ctx, request)
}

func (s serverTx) CreateInvoice(ctx context.Context, request *paymentspb.CreateInvoiceRequest) (resp *paymentspb.CreateInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CreateInvoice(ctx, request)
}

func (s serverTx) AdjustInvoice(ctx context.Context, request *paymentspb.AdjustInvoiceRequest) (resp *paymentspb.AdjustInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.AdjustInvoice(ctx, request)
}

func (s serverTx) PayInvoice(ctx context.Context, request *paymentspb.PayInvoiceRequest) (resp *paymentspb.PayInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.PayInvoice(ctx, request)
}

func (s serverTx) CancelInvoice(ctx context.Context, request *paymentspb.CancelInvoiceRequest) (resp *paymentspb.CancelInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CancelInvoice(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
