package domain

import (
	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
)

const ProductAggregate = "stores.Product"

var (
	ErrProductNameIsBlank     = errors.Wrap(errors.ErrBadRequest, "the product name cannot be blank")
	ErrProductPriceIsNegative = errors.Wrap(errors.ErrBadRequest, "the product price cannot be negative")
	ErrNotAPriceIncrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotAPriceDecrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
)

type Product struct {
	es.Aggregate
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Product)(nil)

func NewProduct(id string) *Product {
	return &Product{
		Aggregate: es.NewAggregate(id, ProductAggregate),
	}
}

func (p *Product) InitProduct(id, storeID, name, description, sku string, price float64) (ddd.Event, error) {
	if name == "" {
		return nil, ErrProductNameIsBlank
	}

	if price < 0 {
		return nil, ErrProductPriceIsNegative
	}

	p.AddEvent(ProductAddedEvent, &ProductAdded{
		StoreID:     storeID,
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
	})

	return ddd.NewEvent(ProductAddedEvent, p), nil
}

// Key implements registry.Registerable
func (Product) Key() string { return ProductAggregate }

func (p *Product) Rebrand(name, description string) (ddd.Event, error) {
	p.AddEvent(ProductRebrandedEvent, &ProductRebranded{
		Name:        name,
		Description: description,
	})

	return ddd.NewEvent(ProductRebrandedEvent, p), nil
}

func (p *Product) IncreasePrice(price float64) (ddd.Event, error) {
	if price < p.Price {
		return nil, ErrNotAPriceIncrease
	}

	delta := price - p.Price
	p.AddEvent(ProductPriceIncreasedEvent, &ProductPriceChanged{
		Delta: delta,
	})

	return ddd.NewEvent(ProductPriceIncreasedEvent, ProductPriceDelta{
		Product: p,
		Delta:   delta,
	}), nil
}

func (p *Product) DecreasePrice(price float64) (ddd.Event, error) {
	if price > p.Price {
		return nil, ErrNotAPriceDecrease
	}

	delta := price - p.Price
	p.AddEvent(ProductPriceDecreasedEvent, &ProductPriceChanged{
		Delta: delta,
	})

	return ddd.NewEvent(ProductPriceDecreasedEvent, ProductPriceDelta{
		Product: p,
		Delta:   delta,
	}), nil
}

func (p *Product) Remove() (ddd.Event, error) {
	p.AddEvent(ProductRemovedEvent, &ProductRemoved{})

	return ddd.NewEvent(ProductRemovedEvent, p), nil
}

func (p *Product) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ProductAdded:
		p.StoreID = payload.StoreID
		p.Name = payload.Name
		p.Description = payload.Description
		p.SKU = payload.SKU
		p.Price = payload.Price

	case *ProductRebranded:
		p.Name = payload.Name
		p.Description = payload.Description

	case *ProductPriceChanged:
		p.Price = p.Price + payload.Delta

	case *ProductRemoved:
		// noop

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", p, event.EventName(), payload)
	}

	return nil
}

func (p *Product) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ProductV1:
		p.StoreID = ss.StoreID
		p.Name = ss.Name
		p.Description = ss.Description
		p.SKU = ss.SKU
		p.Price = ss.Price

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", p, snapshot)
	}

	return nil
}

func (p Product) ToSnapshot() es.Snapshot {
	return ProductV1{
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Price:       p.Price,
	}
}
