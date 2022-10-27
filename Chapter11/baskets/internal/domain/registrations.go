package domain

import (
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/serdes"
)

const (
	BasketStartedEvent     = "baskets.BasketStarted"
	BasketItemAddedEvent   = "baskets.BasketItemAdded"
	BasketItemRemovedEvent = "baskets.BasketItemRemoved"
	BasketCanceledEvent    = "baskets.BasketCanceled"
	BasketCheckedOutEvent  = "baskets.BasketCheckedOut"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// Basket
	if err := serde.Register(Basket{}, func(v interface{}) error {
		basket := v.(*Basket)
		basket.Aggregate = es.NewAggregate("", BasketAggregate)
		basket.Items = make(map[string]Item)
		return nil
	}); err != nil {
		return err
	}
	// basket events
	if err := serde.Register(BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(BasketCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(BasketCheckedOut{}); err != nil {
		return err
	}
	if err := serde.Register(BasketItemAdded{}); err != nil {
		return err
	}
	if err := serde.Register(BasketItemRemoved{}); err != nil {
		return err
	}
	// basket snapshots
	if err := serde.RegisterKey(BasketV1{}.SnapshotName(), BasketV1{}); err != nil {
		return err
	}

	return nil
}

func (Basket) Key() string { return BasketAggregate }

func (BasketStarted) Key() string     { return BasketStartedEvent }
func (BasketItemAdded) Key() string   { return BasketItemAddedEvent }
func (BasketItemRemoved) Key() string { return BasketItemRemovedEvent }
func (BasketCanceled) Key() string    { return BasketCanceledEvent }
func (BasketCheckedOut) Key() string  { return BasketCheckedOutEvent }
