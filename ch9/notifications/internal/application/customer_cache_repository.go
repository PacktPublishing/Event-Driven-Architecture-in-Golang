package application

import (
	"context"
)

type CustomerCacheRepository interface {
	Add(ctx context.Context, customerID, name, smsNumber string) error
	UpdateSmsNumber(ctx context.Context, customerID, smsNumber string) error
	CustomerRepository
}
