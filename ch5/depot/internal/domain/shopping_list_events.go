package domain

const (
	ShoppingListCreatedEvent   = "depot.ShoppingListCreated"
	ShoppingListCanceledEvent  = "depot.ShoppingListCanceled"
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

type ShoppingListAssigned struct {
	ShoppingList *ShoppingList
	BotID        string
}

func (ShoppingListAssigned) Key() string { return ShoppingListAssignedEvent }

type ShoppingListCompleted struct {
	ShoppingList *ShoppingList
}

func (ShoppingListCompleted) Key() string { return ShoppingListCompletedEvent }
