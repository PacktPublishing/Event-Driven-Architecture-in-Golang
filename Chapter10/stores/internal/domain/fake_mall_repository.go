package domain

import (
	"context"
)

type FakeMallRepository struct {
	stores map[string]*MallStore
}

var _ MallRepository = (*FakeMallRepository)(nil)

func NewFakeMallRepository() *FakeMallRepository {
	return &FakeMallRepository{
		stores: map[string]*MallStore{},
	}
}

func (r *FakeMallRepository) AddStore(ctx context.Context, storeID, name, location string) error {
	// TODO implement me
	panic("implement me")
}

func (r *FakeMallRepository) SetStoreParticipation(ctx context.Context, storeID string, participating bool) error {
	// TODO implement me
	panic("implement me")
}

func (r *FakeMallRepository) RenameStore(ctx context.Context, storeID, name string) error {
	// TODO implement me
	panic("implement me")
}

func (r *FakeMallRepository) Find(ctx context.Context, storeID string) (*MallStore, error) {
	// TODO implement me
	panic("implement me")
}

func (r *FakeMallRepository) All(ctx context.Context) ([]*MallStore, error) {
	// TODO implement me
	panic("implement me")
}

func (r *FakeMallRepository) AllParticipating(ctx context.Context) ([]*MallStore, error) {
	// TODO implement me
	panic("implement me")
}
