package domain

import (
	"context"
)

type CustomerRepository interface {
	Authorize(ctx context.Context, customerID string) error
}
