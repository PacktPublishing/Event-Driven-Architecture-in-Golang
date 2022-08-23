package commands

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type RebrandStore struct {
	ID   string
	Name string
}

type RebrandStoreHandler struct {
	stores domain.StoreRepository
}

func NewRebrandStoreHandler(stores domain.StoreRepository) RebrandStoreHandler {
	return RebrandStoreHandler{
		stores: stores,
	}
}

func (h RebrandStoreHandler) RebrandStore(ctx context.Context, cmd RebrandStore) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = store.Rebrand(cmd.Name); err != nil {
		return err
	}

	return h.stores.Save(ctx, store)
}
