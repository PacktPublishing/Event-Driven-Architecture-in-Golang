package domain

const (
	BasketStartedEvent     = "baskets.BasketStarted"
	BasketItemAddedEvent   = "baskets.BasketItemAdded"
	BasketItemRemovedEvent = "baskets.BasketItemRemoved"
	BasketCanceledEvent    = "baskets.BasketCanceled"
	BasketCheckedOutEvent  = "baskets.BasketCheckedOut"
)

type BasketStarted struct {
	CustomerID string
}

// Key implements registry.Registerable
func (BasketStarted) Key() string { return BasketStartedEvent }

type BasketItemAdded struct {
	Item Item
}

// Key implements registry.Registerable
func (BasketItemAdded) Key() string { return BasketItemAddedEvent }

type BasketItemRemoved struct {
	ProductID string
	Quantity  int
}

// Key implements registry.Registerable
func (BasketItemRemoved) Key() string { return BasketItemRemovedEvent }

type BasketCanceled struct{}

// Key implements registry.Registerable
func (BasketCanceled) Key() string { return BasketCanceledEvent }

type BasketCheckedOut struct {
	PaymentID  string
	CustomerID string
	Items      map[string]Item
}

// Key implements registry.Registerable
func (BasketCheckedOut) Key() string { return BasketCheckedOutEvent }
