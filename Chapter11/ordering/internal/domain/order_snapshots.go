package domain

type OrderV1 struct {
	CustomerID string
	PaymentID  string
	InvoiceID  string
	ShoppingID string
	Items      []Item
	Status     OrderStatus
}

func (OrderV1) SnapshotName() string { return "ordering.OrderV1" }
