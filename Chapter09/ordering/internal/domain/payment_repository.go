package domain

import (
	"context"
)

type PaymentRepository interface {
	Confirm(ctx context.Context, paymentID string) error
}
