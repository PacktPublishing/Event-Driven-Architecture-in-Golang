package domain

type BasketStatus string

const (
	BasketUnknown      BasketStatus = ""
	BasketIsOpen       BasketStatus = "open"
	BasketIsCanceled   BasketStatus = "canceled"
	BasketIsCheckedOut BasketStatus = "checked_out"
)

func (s BasketStatus) String() string {
	return string(s)
}

func ToBasketStatus(status string) BasketStatus {
	switch status {
	case BasketIsOpen.String():
		return BasketIsOpen
	case BasketIsCanceled.String():
		return BasketIsCanceled
	case BasketIsCheckedOut.String():
		return BasketIsCheckedOut
	default:
		return BasketUnknown
	}
}
