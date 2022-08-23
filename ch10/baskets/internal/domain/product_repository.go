package domain

import (
	"context"
)

type ProductRepository interface {
	Find(ctx context.Context, productID string) (*Product, error)
}
