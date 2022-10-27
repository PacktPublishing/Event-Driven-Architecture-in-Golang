package models

type InvoiceStatus string

const (
	InvoiceUnknown  InvoiceStatus = ""
	InvoicePending  InvoiceStatus = "pending"
	InvoicePaid     InvoiceStatus = "paid"
	InvoiceCanceled InvoiceStatus = "canceled"
)

type Invoice struct {
	ID      string
	OrderID string
	Amount  float64
	Status  InvoiceStatus
}

func (s InvoiceStatus) String() string {
	switch s {
	case InvoicePending, InvoicePaid, InvoiceCanceled:
		return string(s)
	default:
		return ""
	}
}
