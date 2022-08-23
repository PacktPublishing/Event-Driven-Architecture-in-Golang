package domain

import (
	"context"
)

type FakeStoreRepository struct {
	stores map[string]*Store
}

func NewFakeStoreRepository() *FakeStoreRepository {
	return &FakeStoreRepository{stores: map[string]*Store{}}
}

var _ StoreRepository = (*FakeStoreRepository)(nil)

func (r *FakeStoreRepository) Load(ctx context.Context, storeID string) (*Store, error) {
	if store, exists := r.stores[storeID]; exists {
		return store, nil
	}

	return NewStore(storeID), nil
}

func (r *FakeStoreRepository) Save(ctx context.Context, store *Store) error {
	for _, event := range store.Events() {
		if err := store.ApplyEvent(event); err != nil {
			return err
		}
	}

	r.stores[store.ID()] = store

	return nil
}

func (r *FakeStoreRepository) Reset(stores ...*Store) {
	r.stores = make(map[string]*Store)

	for _, store := range stores {
		r.stores[store.ID()] = store
	}
}
