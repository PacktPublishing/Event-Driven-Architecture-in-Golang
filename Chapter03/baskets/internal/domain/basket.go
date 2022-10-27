package domain

import (
	"sort"

	"github.com/stackus/errors"
)

var (
	ErrBasketHasNoItems         = errors.Wrap(errors.ErrBadRequest, "the basket has no items")
	ErrBasketCannotBeModified   = errors.Wrap(errors.ErrBadRequest, "the basket cannot be modified")
	ErrBasketCannotBeCancelled  = errors.Wrap(errors.ErrBadRequest, "the basket cannot be cancelled")
	ErrQuantityCannotBeNegative = errors.Wrap(errors.ErrBadRequest, "the item quantity cannot be negative")
	ErrBasketIDCannotBeBlank    = errors.Wrap(errors.ErrBadRequest, "the basket id cannot be blank")
	ErrPaymentIDCannotBeBlank   = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
	ErrCustomerIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
)

type BasketStatus string

const (
	BasketUnknown    BasketStatus = ""
	BasketOpen       BasketStatus = "open"
	BasketCancelled  BasketStatus = "cancelled"
	BasketCheckedOut BasketStatus = "checked_out"
)

func (s BasketStatus) String() string {
	switch s {
	case BasketOpen, BasketCancelled, BasketCheckedOut:
		return string(s)
	default:
		return ""
	}
}

type Basket struct {
	ID         string
	CustomerID string
	PaymentID  string
	Items      []Item
	Status     BasketStatus
}

func StartBasket(id, customerID string) (*Basket, error) {
	if id == "" {
		return nil, ErrBasketIDCannotBeBlank
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	basket := &Basket{
		ID:         id,
		CustomerID: customerID,
		Status:     BasketOpen,
		Items:      []Item{},
	}

	return basket, nil
}

func (b Basket) IsCancellable() bool {
	return b.Status == BasketOpen
}

func (b Basket) IsOpen() bool {
	return b.Status == BasketOpen
}

func (b *Basket) Cancel() error {
	if !b.IsCancellable() {
		return ErrBasketCannotBeCancelled
	}

	b.Status = BasketCancelled
	b.Items = []Item{}

	return nil
}

func (b *Basket) Checkout(paymentID string) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if len(b.Items) == 0 {
		return ErrBasketHasNoItems
	}

	if paymentID == "" {
		return ErrPaymentIDCannotBeBlank
	}

	b.PaymentID = paymentID
	b.Status = BasketCheckedOut

	return nil
}

func (b *Basket) AddItem(store *Store, product *Product, quantity int) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if quantity < 0 {
		return ErrQuantityCannotBeNegative
	}

	for i, item := range b.Items {
		if item.ProductID == product.ID && item.StoreID == product.StoreID {
			b.Items[i].Quantity += quantity
			return nil
		}
	}

	b.Items = append(b.Items, Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     quantity,
	})

	sort.Slice(b.Items, func(i, j int) bool {
		return b.Items[i].StoreName <= b.Items[j].StoreName && b.Items[i].ProductName < b.Items[j].ProductName
	})

	return nil
}

func (b *Basket) RemoveItem(product *Product, quantity int) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if quantity < 0 {
		return ErrQuantityCannotBeNegative
	}

	for i, item := range b.Items {
		if item.ProductID == product.ID && item.StoreID == product.StoreID {
			b.Items[i].Quantity -= quantity

			if b.Items[i].Quantity < 1 {
				b.Items = append(b.Items[:i], b.Items[i+1:]...)
			}
			return nil
		}
	}

	return nil
}
