package domain

const (
	ShoppingListCreatedEvent   = "depot.ShoppingListCreated"
	ShoppingListCanceledEvent  = "depot.ShoppingListCanceled"
	ShoppingListInitiatedEvent = "depot.ShoppingListInitiated"
	ShoppingListAssignedEvent  = "depot.ShoppingListAssigned"
	ShoppingListCompletedEvent = "depot.ShoppingListCompleted"
)

type ShoppingListCreated struct {
	ShoppingList *ShoppingList
}

func (ShoppingListCreated) Key() string { return ShoppingListCreatedEvent }

type ShoppingListCanceled struct {
	ShoppingList *ShoppingList
}

func (ShoppingListCanceled) Key() string { return ShoppingListCanceledEvent }

type ShoppingListInitiated struct {
	ShoppingList *ShoppingList
}

func (ShoppingListInitiated) Key() string { return ShoppingListInitiatedEvent }

type ShoppingListAssigned struct {
	ShoppingList *ShoppingList
	BotID        string
}

func (ShoppingListAssigned) Key() string { return ShoppingListAssignedEvent }

type ShoppingListCompleted struct {
	ShoppingList *ShoppingList
}

func (ShoppingListCompleted) Key() string { return ShoppingListCompletedEvent }
