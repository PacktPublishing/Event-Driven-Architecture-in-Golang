package application

import (
	"context"

	"eda-in-golang/search/internal/models"
)

type ProductRepository interface {
	Find(ctx context.Context, productID string) (*models.Product, error)
}
