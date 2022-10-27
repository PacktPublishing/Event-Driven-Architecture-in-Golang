package domain

import (
	"sort"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
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

type Basket struct {
	ddd.AggregateBase
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
		AggregateBase: ddd.AggregateBase{
			ID: id,
		},
		CustomerID: customerID,
		Status:     BasketIsOpen,
		Items:      []Item{},
	}

	basket.AddEvent(&BasketStarted{
		Basket: basket,
	})

	return basket, nil
}

func (b Basket) IsCancellable() bool {
	return b.Status == BasketIsOpen
}

func (b Basket) IsOpen() bool {
	return b.Status == BasketIsOpen
}

func (b *Basket) Cancel() error {
	if !b.IsCancellable() {
		return ErrBasketCannotBeCancelled
	}

	b.Status = BasketIsCanceled
	b.Items = []Item{}

	b.AddEvent(&BasketCanceled{
		Basket: b,
	})

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
	b.Status = BasketIsCheckedOut

	b.AddEvent(&BasketCheckedOut{
		Basket: b,
	})

	return nil
}

func (b *Basket) hasProduct(product *Product) (int, bool) {
	for i, item := range b.Items {
		if item.ProductID == product.ID && item.StoreID == product.StoreID {
			return i, true
		}
	}

	return -1, false
}

func (b *Basket) AddItem(store *Store, product *Product, quantity int) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if quantity < 0 {
		return ErrQuantityCannotBeNegative
	}

	item := Item{
		StoreID:      store.ID,
		ProductID:    product.ID,
		StoreName:    store.Name,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     quantity,
	}

	if i, exists := b.hasProduct(product); exists {
		b.Items[i].Quantity += quantity
	} else {
		b.Items = append(b.Items, item)

		sort.Slice(b.Items, func(i, j int) bool {
			return b.Items[i].StoreName <= b.Items[j].StoreName && b.Items[i].ProductName < b.Items[j].ProductName
		})
	}

	b.AddEvent(&BasketItemAdded{
		Basket: b,
		Item:   item,
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

	if i, exists := b.hasProduct(product); exists {
		b.Items[i].Quantity -= quantity

		item := b.Items[i]
		item.Quantity = quantity

		if b.Items[i].Quantity < 1 {
			b.Items = append(b.Items[:i], b.Items[i+1:]...)
		}

		b.AddEvent(&BasketItemRemoved{
			Basket: b,
			Item:   item,
		})
	}

	return nil
}
