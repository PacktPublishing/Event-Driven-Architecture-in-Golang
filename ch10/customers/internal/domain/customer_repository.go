package domain

import (
	"context"
)

type CustomerRepository interface {
	Save(ctx context.Context, customer *Customer) error
	Find(ctx context.Context, customerID string) (*Customer, error)
	Update(ctx context.Context, customer *Customer) error
}
