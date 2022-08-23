package domain

const (
	ProductAddedEvent          = "stores.ProductAdded"
	ProductRebrandedEvent      = "stores.ProductRebranded"
	ProductPriceIncreasedEvent = "stores.ProductPriceIncreased"
	ProductPriceDecreasedEvent = "stores.ProductPriceDecreased"
	ProductRemovedEvent        = "stores.ProductRemoved"
)

type ProductAdded struct {
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

// Key implements registry.Registerable
func (ProductAdded) Key() string { return ProductAddedEvent }

type ProductRebranded struct {
	Name        string
	Description string
}

// Key implements registry.Registerable
func (ProductRebranded) Key() string { return ProductRebrandedEvent }

type ProductPriceChanged struct {
	Delta float64
}

type ProductRemoved struct{}

// Key implements registry.Registerable
func (ProductRemoved) Key() string { return ProductRemovedEvent }
