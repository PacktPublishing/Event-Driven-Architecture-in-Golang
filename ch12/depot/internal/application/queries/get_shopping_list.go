package queries

import (
	"context"

	"eda-in-golang/depot/internal/domain"
)

type GetShoppingList struct {
	ID string
}

type GetShoppingListHandler struct {
	shoppingLists domain.ShoppingListRepository
}

func NewGetShoppingListHandler(shoppingLists domain.ShoppingListRepository) GetShoppingListHandler {
	return GetShoppingListHandler{shoppingLists: shoppingLists}
}

func (h GetShoppingListHandler) GetShoppingList(ctx context.Context, query GetShoppingList) (*domain.ShoppingList, error) {

	return h.shoppingLists.Find(ctx, query.ID)
}
