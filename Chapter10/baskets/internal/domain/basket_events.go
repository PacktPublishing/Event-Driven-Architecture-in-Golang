package domain

type BasketStarted struct {
	CustomerID string
}

type BasketItemAdded struct {
	Item Item
}

type BasketItemRemoved struct {
	ProductID string
	Quantity  int
}

type BasketCanceled struct{}

type BasketCheckedOut struct {
	PaymentID string
}
