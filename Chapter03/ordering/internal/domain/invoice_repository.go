package domain

import (
	"context"
)

type InvoiceRepository interface {
	Save(ctx context.Context, orderID, paymentID string, amount float64) error
	Delete(ctx context.Context, invoiceID string) error
}
