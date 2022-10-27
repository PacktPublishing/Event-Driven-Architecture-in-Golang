package models

import (
	"time"
)

type Order struct {
	OrderID      string
	CustomerID   string
	CustomerName string
	Items        []Item
	Total        float64
	Status       string
	CreatedAt    time.Time
}

type Item struct {
	ProductID   string
	StoreID     string
	ProductName string
	StoreName   string
	Price       float64
	Quantity    int
}
