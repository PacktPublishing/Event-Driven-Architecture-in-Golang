package application

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/ddd"
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
		baskets   domain.BasketRepository
		stores    domain.StoreRepository
		products  domain.ProductRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

var _ App = (*Application)(nil)

func New(baskets domain.BasketRepository, stores domain.StoreRepository, products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		baskets:   baskets,
		stores:    stores,
		products:  products,
		publisher: publisher,
	}
}

func (a Application) StartBasket(ctx context.Context, start StartBasket) error {
	basket, err := a.baskets.Load(ctx, start.ID)
	if err != nil {
		return err
	}

	event, err := basket.Start(start.CustomerID)
	if err != nil {
		return err
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		return err
	}

	return a.publisher.Publish(ctx, event)
}

func (a Application) CancelBasket(ctx context.Context, cancel CancelBasket) error {
	basket, err := a.baskets.Load(ctx, cancel.ID)
	if err != nil {
		return err
	}

	event, err := basket.Cancel()
	if err != nil {
		return err
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		return err
	}

	return a.publisher.Publish(ctx, event)
}

func (a Application) CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error {
	basket, err := a.baskets.Load(ctx, checkout.ID)
	if err != nil {
		return err
	}

	event, err := basket.Checkout(checkout.PaymentID)
	if err != nil {
		return errors.Wrap(err, "baskets checkout")
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		return errors.Wrap(err, "basket checkout")
	}

	return a.publisher.Publish(ctx, event)
}

func (a Application) AddItem(ctx context.Context, add AddItem) error {
	basket, err := a.baskets.Load(ctx, add.ID)
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

	if err = a.baskets.Save(ctx, basket); err != nil {
		return err
	}

	return nil
}

func (a Application) RemoveItem(ctx context.Context, remove RemoveItem) error {
	product, err := a.products.Find(ctx, remove.ProductID)
	if err != nil {
		return err
	}

	basket, err := a.baskets.Load(ctx, remove.ID)
	if err != nil {
		return err
	}

	err = basket.RemoveItem(product, remove.Quantity)
	if err != nil {
		return err
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		return err
	}

	return nil
}

func (a Application) GetBasket(ctx context.Context, get GetBasket) (*domain.Basket, error) {
	return a.baskets.Load(ctx, get.ID)
}
