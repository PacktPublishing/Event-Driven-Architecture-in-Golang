package domain

const (
	OrderCreatedEvent   = "ordering.OrderCreated"
	OrderRejectedEvent  = "ordering.OrderRejected"
	OrderApprovedEvent  = "ordering.OrderApproved"
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

type OrderRejected struct{}

func (OrderRejected) Key() string { return OrderRejectedEvent }

type OrderApproved struct {
	ShoppingID string
}

func (OrderApproved) Key() string { return OrderApprovedEvent }

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
