package commands

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

type DisableParticipation struct {
	ID string
}

type DisableParticipationHandler struct {
	stores domain.StoreRepository
}

func NewDisableParticipationHandler(stores domain.StoreRepository) DisableParticipationHandler {
	return DisableParticipationHandler{
		stores: stores,
	}
}

func (h DisableParticipationHandler) DisableParticipation(ctx context.Context, cmd DisableParticipation) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = store.DisableParticipation(); err != nil {
		return err
	}

	return h.stores.Save(ctx, store)
}
