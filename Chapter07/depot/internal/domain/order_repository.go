package domain

import (
	"context"
)

type OrderRepository interface {
	Ready(ctx context.Context, orderID string) error
}
