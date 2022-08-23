package domain

import (
	"context"
)

type OrderRepository interface {
	Save(ctx context.Context, basket *Basket) (string, error)
}
