package models

type InvoiceStatus string

const (
	InvoiceIsUnknown  InvoiceStatus = ""
	InvoiceIsPending  InvoiceStatus = "pending"
	InvoiceIsPaid     InvoiceStatus = "paid"
	InvoiceIsCanceled InvoiceStatus = "canceled"
)

type Invoice struct {
	ID      string
	OrderID string
	Amount  float64
	Status  InvoiceStatus
}

func (s InvoiceStatus) String() string {
	switch s {
	case InvoiceIsPending, InvoiceIsPaid, InvoiceIsCanceled:
		return string(s)
	default:
		return ""
	}
}
