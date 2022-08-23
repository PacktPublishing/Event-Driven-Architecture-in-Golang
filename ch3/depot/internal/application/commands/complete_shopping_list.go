package commands

import (
	"context"

	"eda-in-golang/depot/internal/domain"
)

type CompleteShoppingList struct {
	ID string
}

type CompleteShoppingListHandler struct {
	shoppingLists domain.ShoppingListRepository
	orders        domain.OrderRepository
}

func NewCompleteShoppingListHandler(shoppingLists domain.ShoppingListRepository, orders domain.OrderRepository,
) CompleteShoppingListHandler {
	return CompleteShoppingListHandler{
		shoppingLists: shoppingLists,
		orders:        orders,
	}
}

func (h CompleteShoppingListHandler) CompleteShoppingList(ctx context.Context, cmd CompleteShoppingList) error {
	list, err := h.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = list.Complete()
	if err != nil {
		return err
	}

	err = h.orders.Ready(ctx, list.OrderID)
	if err != nil {
		return err
	}

	return h.shoppingLists.Update(ctx, list)
}
