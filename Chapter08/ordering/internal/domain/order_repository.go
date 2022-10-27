package domain

import (
	"context"
)

type OrderRepository interface {
	Load(ctx context.Context, orderID string) (*Order, error)
	Save(ctx context.Context, order *Order) error
}
