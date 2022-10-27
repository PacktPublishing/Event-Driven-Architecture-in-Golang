package domain

import (
	"context"
)

type ProductRepository interface {
	Load(ctx context.Context, id string) (*Product, error)
	Save(ctx context.Context, product *Product) error
}
