package application

import (
	"context"

	"eda-in-golang/search/internal/models"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}
