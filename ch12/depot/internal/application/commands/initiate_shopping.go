package commands

import (
	"context"

	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/ddd"
)

type InitiateShopping struct {
	ID string
}

type InitiateShoppingHandler struct {
	shoppingLists   domain.ShoppingListRepository
	domainPublisher ddd.EventPublisher[ddd.AggregateEvent]
}

func NewInitiateShoppingHandler(lists domain.ShoppingListRepository, publisher ddd.EventPublisher[ddd.AggregateEvent]) InitiateShoppingHandler {
	return InitiateShoppingHandler{
		shoppingLists:   lists,
		domainPublisher: publisher,
	}
}

func (h InitiateShoppingHandler) InitiateShopping(ctx context.Context, cmd InitiateShopping) error {
	list, err := h.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = list.Initiate()
	if err != nil {
		return err
	}

	if err = h.shoppingLists.Update(ctx, list); err != nil {
		return err
	}

	// publish domain events
	if err = h.domainPublisher.Publish(ctx, list.Events()...); err != nil {
		return err
	}

	return nil
}
