package domain

import (
	"context"
)

type OrderRepository interface {
	Save(ctx context.Context, paymentID, customerID string, basketItems map[string]Item) (string, error)
}
