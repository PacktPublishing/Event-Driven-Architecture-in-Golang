package domain

type BasketV1 struct {
	CustomerID string
	PaymentID  string
	Items      map[string]Item
	Status     BasketStatus
}

func (BasketV1) SnapshotName() string { return "baskets.BasketV1" }
