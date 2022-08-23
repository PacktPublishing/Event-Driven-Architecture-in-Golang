package domain

type OrderCreated struct {
	Order *Order
}

func (OrderCreated) EventName() string { return "ordering.OrderCreated" }

type OrderCanceled struct {
	Order *Order
}

func (OrderCanceled) EventName() string { return "ordering.OrderCanceled" }

type OrderReadied struct {
	Order *Order
}

func (OrderReadied) EventName() string { return "ordering.OrderReadied" }

type OrderCompleted struct {
	Order *Order
}

func (OrderCompleted) EventName() string { return "ordering.OrderCompleted" }
