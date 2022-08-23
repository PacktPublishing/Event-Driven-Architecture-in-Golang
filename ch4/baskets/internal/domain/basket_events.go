package domain

type BasketStarted struct {
	Basket *Basket
}

func (BasketStarted) EventName() string { return "baskets.BasketStarted" }

type BasketItemAdded struct {
	Basket *Basket
	Item   Item
}

func (BasketItemAdded) EventName() string { return "baskets.BasketItemAdded" }

type BasketItemRemoved struct {
	Basket *Basket
	Item   Item
}

func (BasketItemRemoved) EventName() string { return "baskets.BasketItemRemoved" }

type BasketCanceled struct {
	Basket *Basket
}

func (BasketCanceled) EventName() string { return "baskets.BasketCanceled" }

type BasketCheckedOut struct {
	Basket *Basket
}

func (BasketCheckedOut) EventName() string { return "baskets.BasketCheckedOut" }
