package commands

import (
	"context"

	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/ddd"
)

type CompleteShoppingList struct {
	ID string
}

type CompleteShoppingListHandler struct {
	shoppingLists   domain.ShoppingListRepository
	domainPublisher ddd.EventPublisher
}

func NewCompleteShoppingListHandler(shoppingLists domain.ShoppingListRepository, domainPublisher ddd.EventPublisher) CompleteShoppingListHandler {
	return CompleteShoppingListHandler{
		shoppingLists:   shoppingLists,
		domainPublisher: domainPublisher,
	}
}

func (h CompleteShoppingListHandler) CompleteShoppingList(ctx context.Context, cmd CompleteShoppingList) error {
	list, err := h.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = list.Complete(); err != nil {
		return err
	}

	if err = h.shoppingLists.Update(ctx, list); err != nil {
		return nil
	}

	// publish domain events
	if err = h.domainPublisher.Publish(ctx, list.GetEvents()...); err != nil {
		return err
	}

	return nil
}
