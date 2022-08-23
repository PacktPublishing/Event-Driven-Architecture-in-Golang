package domain

import (
	"context"
)

type MallStore struct {
	ID            string
	Name          string
	Location      string
	Participating bool
}

type MallRepository interface {
	AddStore(ctx context.Context, storeID, name, location string) error
	SetStoreParticipation(ctx context.Context, storeID string, participating bool) error
	RenameStore(ctx context.Context, storeID, name string) error
	Find(ctx context.Context, storeID string) (*MallStore, error)
	All(ctx context.Context) ([]*MallStore, error)
	AllParticipating(ctx context.Context) ([]*MallStore, error)
}
