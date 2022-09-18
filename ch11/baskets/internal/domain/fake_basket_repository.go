package domain

import (
	"context"
)

type FakeBasketRepository struct {
	baskets map[string]*Basket
}

func NewFakeBasketRepository() *FakeBasketRepository {
	return &FakeBasketRepository{baskets: map[string]*Basket{}}
}

var _ BasketRepository = (*FakeBasketRepository)(nil)

func (r *FakeBasketRepository) Load(ctx context.Context, basketID string) (*Basket, error) {
	if basket, exists := r.baskets[basketID]; exists {
		return basket, nil
	}

	return NewBasket(basketID), nil
}

func (r *FakeBasketRepository) Save(ctx context.Context, basket *Basket) error {
	r.baskets[basket.ID()] = basket

	return nil
}

func (r *FakeBasketRepository) Reset(baskets ...*Basket) {
	r.baskets = make(map[string]*Basket)

	for _, basket := range baskets {
		r.baskets[basket.ID()] = basket
	}
}
