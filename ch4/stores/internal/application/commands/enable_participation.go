package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

type EnableParticipation struct {
	ID string
}

type EnableParticipationHandler struct {
	stores          domain.StoreRepository
	domainPublisher ddd.EventPublisher
}

func NewEnableParticipationHandler(stores domain.StoreRepository, domainPublisher ddd.EventPublisher) EnableParticipationHandler {
	return EnableParticipationHandler{
		stores:          stores,
		domainPublisher: domainPublisher,
	}
}

func (h EnableParticipationHandler) EnableParticipation(ctx context.Context, cmd EnableParticipation) error {
	store, err := h.stores.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = store.EnableParticipation(); err != nil {
		return err
	}

	if err = h.stores.Update(ctx, store); err != nil {
		return err
	}

	if err = h.domainPublisher.Publish(ctx, store.GetEvents()...); err != nil {
		return err
	}

	return nil
}
