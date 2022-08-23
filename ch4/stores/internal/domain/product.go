package domain

import (
	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
)

var (
	ErrProductNameIsBlank     = errors.Wrap(errors.ErrBadRequest, "the product name cannot be blank")
	ErrProductPriceIsNegative = errors.Wrap(errors.ErrBadRequest, "the product price cannot be negative")
)

type Product struct {
	ddd.AggregateBase
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

func CreateProduct(id, storeID, name, description, sku string, price float64) (*Product, error) {
	if name == "" {
		return nil, ErrProductNameIsBlank
	}

	if price < 0 {
		return nil, ErrProductPriceIsNegative
	}

	product := &Product{
		AggregateBase: ddd.AggregateBase{
			ID: id,
		},
		StoreID:     storeID,
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
	}

	product.AddEvent(&ProductAdded{
		Product: product,
	})

	return product, nil
}

func (p *Product) Remove() error {
	p.AddEvent(&ProductRemoved{
		Product: p,
	})

	return nil
}
