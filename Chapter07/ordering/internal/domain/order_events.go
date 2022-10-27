package domain

const (
	OrderCreatedEvent   = "ordering.OrderCreated"
	OrderCanceledEvent  = "ordering.OrderCanceled"
	OrderReadiedEvent   = "ordering.OrderReadied"
	OrderCompletedEvent = "ordering.OrderCompleted"
)

type OrderCreated struct {
	CustomerID string
	PaymentID  string
	ShoppingID string
	Items      []Item
}

func (OrderCreated) Key() string { return OrderCreatedEvent }

type OrderCanceled struct {
	CustomerID string
	PaymentID  string
}

func (OrderCanceled) Key() string { return OrderCanceledEvent }

type OrderReadied struct {
	CustomerID string
	PaymentID  string
	Total      float64
}

func (OrderReadied) Key() string { return OrderReadiedEvent }

type OrderCompleted struct {
	CustomerID string
	InvoiceID  string
}

func (OrderCompleted) Key() string { return OrderCompletedEvent }
