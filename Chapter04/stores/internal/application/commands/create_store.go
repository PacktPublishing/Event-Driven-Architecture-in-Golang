package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

type (
	CreateStore struct {
		ID       string
		Name     string
		Location string
	}

	CreateStoreHandler struct {
		stores          domain.StoreRepository
		domainPublisher ddd.EventPublisher
	}
)

func NewCreateStoreHandler(stores domain.StoreRepository, domainPublisher ddd.EventPublisher) CreateStoreHandler {
	return CreateStoreHandler{
		stores:          stores,
		domainPublisher: domainPublisher,
	}
}

func (h CreateStoreHandler) CreateStore(ctx context.Context, cmd CreateStore) error {
	store, err := domain.CreateStore(cmd.ID, cmd.Name, cmd.Location)
	if err != nil {
		return err
	}

	if err = h.stores.Save(ctx, store); err != nil {
		return err
	}

	if err = h.domainPublisher.Publish(ctx, store.GetEvents()...); err != nil {
		return err
	}

	return nil
}
