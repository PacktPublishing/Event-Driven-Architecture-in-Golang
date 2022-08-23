package domain

type ProductAdded struct {
	Product *Product
}

func (ProductAdded) EventName() string { return "stores.ProductAdded" }

type ProductRemoved struct {
	Product *Product
}

func (ProductRemoved) EventName() string { return "stores.ProductRemoved" }
