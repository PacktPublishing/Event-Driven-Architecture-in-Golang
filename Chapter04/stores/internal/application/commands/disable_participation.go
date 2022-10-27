package commands

import (
	"context"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

type DisableParticipation struct {
	ID string
}

type DisableParticipationHandler struct {
	stores          domain.StoreRepository
	domainPublisher ddd.EventPublisher
}

func NewDisableParticipationHandler(stores domain.StoreRepository, domainPublisher ddd.EventPublisher) DisableParticipationHandler {
	return DisableParticipationHandler{
		stores:          stores,
		domainPublisher: domainPublisher,
	}
}

func (h DisableParticipationHandler) DisableParticipation(ctx context.Context, cmd DisableParticipation) error {
	store, err := h.stores.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = store.DisableParticipation(); err != nil {
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
