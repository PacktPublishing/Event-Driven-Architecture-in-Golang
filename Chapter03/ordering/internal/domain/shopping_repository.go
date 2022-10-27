package domain

import (
	"context"
)

type ShoppingRepository interface {
	Create(ctx context.Context, order *Order) (string, error)
	Cancel(ctx context.Context, shoppingID string) error
}
