package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"eda-in-golang/payments/internal/application"
	"eda-in-golang/payments/paymentspb"
)

type server struct {
	app application.App
	paymentspb.UnimplementedPaymentsServiceServer
}

var _ paymentspb.PaymentsServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	paymentspb.RegisterPaymentsServiceServer(registrar, server{app: app})
	return nil
}

func (s server) AuthorizePayment(ctx context.Context, request *paymentspb.AuthorizePaymentRequest,
) (*paymentspb.AuthorizePaymentResponse, error) {
	id := uuid.New().String()
	err := s.app.AuthorizePayment(ctx, application.AuthorizePayment{
		ID:         id,
		CustomerID: request.GetCustomerId(),
		Amount:     request.GetAmount(),
	})
	return &paymentspb.AuthorizePaymentResponse{Id: id}, err
}

func (s server) ConfirmPayment(ctx context.Context, request *paymentspb.ConfirmPaymentRequest,
) (*paymentspb.ConfirmPaymentResponse, error) {
	err := s.app.ConfirmPayment(ctx, application.ConfirmPayment{
		ID: request.GetId(),
	})
	return &paymentspb.ConfirmPaymentResponse{}, err
}

func (s server) CreateInvoice(ctx context.Context, request *paymentspb.CreateInvoiceRequest,
) (*paymentspb.CreateInvoiceResponse, error) {
	id := uuid.New().String()
	err := s.app.CreateInvoice(ctx, application.CreateInvoice{
		ID:      id,
		OrderID: request.GetOrderId(),
		Amount:  request.GetAmount(),
	})
	return &paymentspb.CreateInvoiceResponse{
		Id: id,
	}, err
}

func (s server) AdjustInvoice(ctx context.Context, request *paymentspb.AdjustInvoiceRequest,
) (*paymentspb.AdjustInvoiceResponse, error) {
	err := s.app.AdjustInvoice(ctx, application.AdjustInvoice{
		ID:     request.GetId(),
		Amount: request.GetAmount(),
	})
	return &paymentspb.AdjustInvoiceResponse{}, err
}

func (s server) PayInvoice(ctx context.Context, request *paymentspb.PayInvoiceRequest) (*paymentspb.PayInvoiceResponse,
	error,
) {
	err := s.app.PayInvoice(ctx, application.PayInvoice{
		ID: request.GetId(),
	})
	return &paymentspb.PayInvoiceResponse{}, err
}

func (s server) CancelInvoice(ctx context.Context, request *paymentspb.CancelInvoiceRequest,
) (*paymentspb.CancelInvoiceResponse, error) {
	err := s.app.CancelInvoice(ctx, application.CancelInvoice{
		ID: request.GetId(),
	})
	return &paymentspb.CancelInvoiceResponse{}, err
}
