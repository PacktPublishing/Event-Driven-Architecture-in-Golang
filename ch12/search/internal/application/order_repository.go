package application

import (
	"context"

	"eda-in-golang/search/internal/models"
)

type OrderRepository interface {
	Add(ctx context.Context, order *models.Order) error
	UpdateStatus(ctx context.Context, orderID, status string) error
	Search(ctx context.Context, search SearchOrders) ([]*models.Order, error)
	Get(ctx context.Context, orderID string) (*models.Order, error)
}
