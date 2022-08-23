package commands

import (
	"context"

	"eda-in-golang/depot/internal/domain"
)

type AssignShoppingList struct {
	ID    string
	BotID string
}

type AssignShoppingListHandler struct {
	shoppingLists domain.ShoppingListRepository
}

func NewAssignShoppingListHandler(shoppingList domain.ShoppingListRepository) AssignShoppingListHandler {
	return AssignShoppingListHandler{shoppingLists: shoppingList}
}

func (h AssignShoppingListHandler) AssignShoppingList(ctx context.Context, cmd AssignShoppingList) error {
	list, err := h.shoppingLists.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = list.Assign(cmd.BotID)
	if err != nil {
		return err
	}

	return h.shoppingLists.Update(ctx, list)
}
