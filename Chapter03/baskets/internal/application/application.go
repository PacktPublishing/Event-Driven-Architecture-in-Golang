package application

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/baskets/internal/domain"
)

type (
	StartBasket struct {
		ID         string
		CustomerID string
	}

	CancelBasket struct {
		ID string
	}

	CheckoutBasket struct {
		ID        string
		PaymentID string
	}

	AddItem struct {
		ID        string
		ProductID string
		Quantity  int
	}

	RemoveItem struct {
		ID        string
		ProductID string
		Quantity  int
	}

	GetBasket struct {
		ID string
	}

	App interface {
		StartBasket(ctx context.Context, start StartBasket) error
		CancelBasket(ctx context.Context, cancel CancelBasket) error
		CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error
		AddItem(ctx context.Context, add AddItem) error
		RemoveItem(ctx context.Context, remove RemoveItem) error
		GetBasket(ctx context.Context, get GetBasket) (*domain.Basket, error)
	}

	Application struct {
		baskets  domain.BasketRepository
		stores   domain.StoreRepository
		products domain.ProductRepository
		orders   domain.OrderRepository
	}
)

var _ App = (*Application)(nil)

func New(baskets domain.BasketRepository, stores domain.StoreRepository, products domain.ProductRepository, orders domain.OrderRepository) *Application {
	return &Application{
		baskets:  baskets,
		stores:   stores,
		products: products,
		orders:   orders,
	}
}

func (a Application) StartBasket(ctx context.Context, start StartBasket) error {
	basket, err := domain.StartBasket(start.ID, start.CustomerID)
	if err != nil {
		return err
	}

	return a.baskets.Save(ctx, basket)
}

func (a Application) CancelBasket(ctx context.Context, cancel CancelBasket) error {
	basket, err := a.baskets.Find(ctx, cancel.ID)
	if err != nil {
		return err
	}

	err = basket.Cancel()
	if err != nil {
		return err
	}

	return a.baskets.Update(ctx, basket)
}

func (a Application) CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error {
	basket, err := a.baskets.Find(ctx, checkout.ID)
	if err != nil {
		return err
	}

	err = basket.Checkout(checkout.PaymentID)
	if err != nil {
		return errors.Wrap(err, "baskets checkout")
	}

	// submit the basket to the order module
	_, err = a.orders.Save(ctx, basket)
	if err != nil {
		return errors.Wrap(err, "baskets checkout")
	}

	return errors.Wrap(a.baskets.Update(ctx, basket), "basket checkout")
}

func (a Application) AddItem(ctx context.Context, add AddItem) error {
	basket, err := a.baskets.Find(ctx, add.ID)
	if err != nil {
		return err
	}

	product, err := a.products.Find(ctx, add.ProductID)
	if err != nil {
		return err
	}

	store, err := a.stores.Find(ctx, product.StoreID)
	if err != nil {
		return nil
	}
	err = basket.AddItem(store, product, add.Quantity)
	if err != nil {
		return err
	}

	return a.baskets.Update(ctx, basket)
}

func (a Application) RemoveItem(ctx context.Context, remove RemoveItem) error {
	product, err := a.products.Find(ctx, remove.ProductID)
	if err != nil {
		return err
	}

	basket, err := a.baskets.Find(ctx, remove.ID)
	if err != nil {
		return err
	}

	err = basket.RemoveItem(product, remove.Quantity)
	if err != nil {
		return err
	}

	return a.baskets.Update(ctx, basket)
}

func (a Application) GetBasket(ctx context.Context, get GetBasket) (*domain.Basket, error) {
	return a.baskets.Find(ctx, get.ID)
}
