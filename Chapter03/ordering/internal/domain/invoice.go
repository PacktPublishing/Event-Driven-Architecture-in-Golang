package domain

type InvoiceID string

// type Invoice struct {
// 	ID     InvoiceID
// 	Amount float64
// }

func (i InvoiceID) String() string {
	return string(i)
}
