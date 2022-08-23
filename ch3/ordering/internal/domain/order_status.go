package domain

type OrderStatus string

const (
	OrderUnknown   OrderStatus = ""
	OrderPending   OrderStatus = "pending"
	OrderInProcess OrderStatus = "in-progress"
	OrderReady     OrderStatus = "ready"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) String() string {
	switch s {
	case OrderPending, OrderInProcess, OrderReady, OrderCompleted, OrderCancelled:
		return string(s)
	default:
		return ""
	}
}

func ToOrderStatus(status string) OrderStatus {
	switch status {
	case OrderPending.String():
		return OrderPending
	case OrderInProcess.String():
		return OrderInProcess
	case OrderReady.String():
		return OrderReady
	case OrderCancelled.String():
		return OrderCancelled
	case OrderCompleted.String():
		return OrderCompleted
	default:
		return OrderUnknown
	}
}
