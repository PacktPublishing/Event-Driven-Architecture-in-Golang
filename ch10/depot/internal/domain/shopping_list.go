package domain

import (
	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
)

const ShoppingListAggregate = "depot.ShoppingList"

var (
	ErrShoppingCannotBeCanceled  = errors.Wrap(errors.ErrBadRequest, "the shopping list cannot be canceled")
	ErrShoppingCannotBeInitiated = errors.Wrap(errors.ErrBadRequest, "the shopping list cannot be initiated")
	ErrShoppingCannotBeAssigned  = errors.Wrap(errors.ErrBadRequest, "the shopping list cannot be assigned")
	ErrShoppingCannotBeCompleted = errors.Wrap(errors.ErrBadRequest, "the shopping list cannot be completed")
)

type ShoppingList struct {
	ddd.Aggregate
	OrderID       string
	Stops         Stops
	AssignedBotID string
	Status        ShoppingListStatus
}

func NewShoppingList(id string) *ShoppingList {
	return &ShoppingList{
		Aggregate: ddd.NewAggregate(id, ShoppingListAggregate),
	}
}

func CreateShoppingList(id, orderID string) *ShoppingList {
	shoppingList := NewShoppingList(id)
	shoppingList.OrderID = orderID
	shoppingList.Status = ShoppingListIsPending
	shoppingList.Stops = make(Stops)

	shoppingList.AddEvent(ShoppingListCreatedEvent, &ShoppingListCreated{
		ShoppingList: shoppingList,
	})

	return shoppingList
}

func (ShoppingList) Key() string { return ShoppingListAggregate }

func (sl *ShoppingList) AddItem(store *Store, product *Product, quantity int) error {
	if _, exists := sl.Stops[store.ID]; !exists {
		sl.Stops[store.ID] = &Stop{
			StoreName:     store.Name,
			StoreLocation: store.Location,
			Items:         make(Items),
		}
	}

	return sl.Stops[store.ID].AddItem(product, quantity)
}

func (sl ShoppingList) isCancelable() bool {
	switch sl.Status {
	case ShoppingListIsPending, ShoppingListIsAvailable, ShoppingListIsAssigned, ShoppingListIsActive:
		return true
	default:
		return false
	}
}

func (sl *ShoppingList) Cancel() error {
	if !sl.isCancelable() {
		return ErrShoppingCannotBeCanceled
	}

	sl.Status = ShoppingListIsCanceled

	sl.AddEvent(ShoppingListCanceledEvent, &ShoppingListCanceled{
		ShoppingList: sl,
	})

	return nil
}

func (sl ShoppingList) isPending() bool {
	return sl.Status == ShoppingListIsPending
}

func (sl *ShoppingList) Initiate() error {
	if !sl.isPending() {
		return ErrShoppingCannotBeInitiated
	}

	sl.AddEvent(ShoppingListInitiatedEvent, &ShoppingListInitiated{
		ShoppingList: sl,
	})

	return nil
}

func (sl ShoppingList) isAssignable() bool {
	return sl.Status == ShoppingListIsAvailable
}

func (sl *ShoppingList) Assign(id string) error {
	if !sl.isAssignable() {
		return ErrShoppingCannotBeAssigned
	}

	sl.AssignedBotID = id
	sl.Status = ShoppingListIsAssigned

	sl.AddEvent(ShoppingListAssignedEvent, &ShoppingListAssigned{
		ShoppingList: sl,
		BotID:        id,
	})

	return nil
}

func (sl ShoppingList) isCompletable() bool {
	return sl.Status == ShoppingListIsAssigned
}

func (sl *ShoppingList) Complete() error {
	if !sl.isCompletable() {
		return ErrShoppingCannotBeCompleted
	}

	sl.Status = ShoppingListIsCompleted

	sl.AddEvent(ShoppingListCompletedEvent, &ShoppingListCompleted{
		ShoppingList: sl,
	})

	return nil
}
