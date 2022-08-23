package commands

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type (
	CreateStore struct {
		ID       string
		Name     string
		Location string
	}

	CreateStoreHandler struct {
		stores domain.StoreRepository
	}
)

func NewCreateStoreHandler(stores domain.StoreRepository) CreateStoreHandler {
	return CreateStoreHandler{
		stores: stores,
	}
}

func (h CreateStoreHandler) CreateStore(ctx context.Context, cmd CreateStore) error {
	store, err := domain.CreateStore(cmd.ID, cmd.Name, cmd.Location)
	if err != nil {
		return err
	}

	return h.stores.Save(ctx, store)
}
