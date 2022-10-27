package domain

import (
	"context"

	"github.com/stackus/errors"
)

type FakeStoreCacheRepository struct {
	stores map[string]*Store
}

var _ StoreCacheRepository = (*FakeStoreCacheRepository)(nil)

func NewFakeStoreCacheRepository() *FakeStoreCacheRepository {
	return &FakeStoreCacheRepository{stores: map[string]*Store{}}
}

func (r *FakeStoreCacheRepository) Add(ctx context.Context, storeID, name, location string) error {
	r.stores[storeID] = &Store{
		ID:       storeID,
		Name:     name,
		Location: location,
	}

	return nil
}

func (r *FakeStoreCacheRepository) Rename(ctx context.Context, storeID, name string) error {
	if store, exists := r.stores[storeID]; exists {
		store.Name = name
	}

	return nil
}

func (r *FakeStoreCacheRepository) Find(ctx context.Context, storeID string) (*Store, error) {
	if store, exists := r.stores[storeID]; exists {
		return store, nil
	}

	return nil, errors.ErrNotFound.Msgf("store with id: `%s` does not exist", storeID)
}

func (r *FakeStoreCacheRepository) Reset(stores ...*Store) {
	r.stores = make(map[string]*Store)

	for _, store := range stores {
		r.stores[store.ID] = store
	}
}
