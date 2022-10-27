package domain

type ProductV1 struct {
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

func (ProductV1) SnapshotName() string { return "stores.ProductV1" }
