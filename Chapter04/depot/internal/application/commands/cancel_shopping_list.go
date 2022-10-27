package commands

import (
	"context"

	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/ddd"
)

type CancelShoppingList struct {
	ID string
}

type CancelShoppingListHandler struct {
	shoppingLists   domain.ShoppingListRepository
	domainPublisher ddd.EventPublisher
}

func NewCancelShoppingListHandler(shoppingLists domain.ShoppingListRepository, domainPublisher ddd.EventPublisher) CancelShoppingListHandler {
	return CancelShoppingListHandler{
		shoppingLists:   shoppingLists,
		domainPublisher: domainPublisher,
	}
}

func (h CancelShoppingListHandler) CancelShoppingList(ctx context.Context, cmd CancelShoppingList) error {
	list, err := h.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = list.Cancel()
	if err != nil {
		return err
	}

	if err = h.shoppingLists.Update(ctx, list); err != nil {
		return err
	}

	// publish domain events
	if err = h.domainPublisher.Publish(ctx, list.GetEvents()...); err != nil {
		return err
	}

	return nil
}
