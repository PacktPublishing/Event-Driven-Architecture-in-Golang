package application

import (
	"context"
	"time"

	"eda-in-golang/search/internal/models"
)

type (
	Filters struct {
		CustomerID string
		After      time.Time
		Before     time.Time
		StoreIDs   []string
		ProductIDs []string
		MinTotal   float64
		MaxTotal   float64
		Status     string
	}
	SearchOrders struct {
		Filters Filters
		Next    string
		Limit   int
	}

	GetOrder struct {
		OrderID string
	}

	Application interface {
		SearchOrders(ctx context.Context, search SearchOrders) ([]*models.Order, error)
		GetOrder(ctx context.Context, get GetOrder) (*models.Order, error)
	}

	app struct {
		orders OrderRepository
	}
)

var _ Application = (*app)(nil)

func New(orders OrderRepository) *app {
	return &app{
		orders: orders,
	}
}

func (a app) SearchOrders(ctx context.Context, search SearchOrders) ([]*models.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (a app) GetOrder(ctx context.Context, get GetOrder) (*models.Order, error) {
	// TODO implement me
	panic("implement me")
}
