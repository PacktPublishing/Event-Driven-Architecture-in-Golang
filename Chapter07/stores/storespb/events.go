package storespb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/serdes"
)

const (
	StoreAggregateChannel = "mallbots.stores.events.Store"

	StoreCreatedEvent              = "storesapi.StoreCreated"
	StoreParticipatingToggledEvent = "storesapi.StoreParticipatingToggled"
	StoreRebrandedEvent            = "storesapi.StoreRebranded"

	ProductAggregateChannel = "mallbots.stores.events.Product"

	ProductAddedEvent          = "storesapi.ProductAdded"
	ProductRebrandedEvent      = "storesapi.ProductRebranded"
	ProductPriceIncreasedEvent = "storesapi.ProductPriceIncreased"
	ProductPriceDecreasedEvent = "storesapi.ProductPriceDecreased"
	ProductRemovedEvent        = "storesapi.ProductRemoved"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Store events
	if err := serde.Register(&StoreCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&StoreParticipationToggled{}); err != nil {
		return err
	}
	if err := serde.Register(&StoreRebranded{}); err != nil {
		return err
	}

	if err := serde.Register(&ProductAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductRebranded{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(ProductPriceIncreasedEvent, &ProductPriceChanged{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(ProductPriceDecreasedEvent, &ProductPriceChanged{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductRemoved{}); err != nil {
		return err
	}

	return nil
}

func (*StoreCreated) Key() string              { return StoreCreatedEvent }
func (*StoreParticipationToggled) Key() string { return StoreParticipatingToggledEvent }
func (*StoreRebranded) Key() string            { return StoreRebrandedEvent }

func (*ProductAdded) Key() string     { return ProductAddedEvent }
func (*ProductRebranded) Key() string { return ProductRebrandedEvent }
func (*ProductRemoved) Key() string   { return ProductRemovedEvent }
