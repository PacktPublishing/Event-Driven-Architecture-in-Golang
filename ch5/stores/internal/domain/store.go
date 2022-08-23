package domain

import (
	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
)

const StoreAggregate = "stores.Store"

var (
	ErrStoreNameIsBlank               = errors.Wrap(errors.ErrBadRequest, "the store name cannot be blank")
	ErrStoreLocationIsBlank           = errors.Wrap(errors.ErrBadRequest, "the store location cannot be blank")
	ErrStoreIsAlreadyParticipating    = errors.Wrap(errors.ErrBadRequest, "the store is already participating")
	ErrStoreIsAlreadyNotParticipating = errors.Wrap(errors.ErrBadRequest, "the store is already not participating")
)

type Store struct {
	es.Aggregate
	Name          string
	Location      string
	Participating bool
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Store)(nil)

func NewStore(id string) *Store {
	return &Store{
		Aggregate: es.NewAggregate(id, StoreAggregate),
	}
}

func CreateStore(id, name, location string) (*Store, error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	if location == "" {
		return nil, ErrStoreLocationIsBlank
	}

	store := NewStore(id)

	store.AddEvent(StoreCreatedEvent, &StoreCreated{
		Name:     name,
		Location: location,
	})

	return store, nil
}

// Key implements registry.Registerable
func (Store) Key() string { return StoreAggregate }

func (s *Store) EnableParticipation() (err error) {
	if s.Participating {
		return ErrStoreIsAlreadyParticipating
	}

	s.AddEvent(StoreParticipationEnabledEvent, &StoreParticipationToggled{
		Participating: true,
	})

	return
}

func (s *Store) DisableParticipation() (err error) {
	if !s.Participating {
		return ErrStoreIsAlreadyNotParticipating
	}

	s.AddEvent(StoreParticipationDisabledEvent, &StoreParticipationToggled{
		Participating: false,
	})

	return
}

func (s *Store) Rebrand(name string) error {
	s.AddEvent(StoreRebrandedEvent, &StoreRebranded{
		Name: name,
	})

	return nil
}

// ApplyEvent implements es.EventApplier
func (s *Store) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *StoreCreated:
		s.Name = payload.Name
		s.Location = payload.Location

	case *StoreParticipationToggled:
		s.Participating = payload.Participating

	case *StoreRebranded:
		s.Name = payload.Name

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", s, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (s *Store) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *StoreV1:
		s.Name = ss.Name
		s.Location = ss.Location
		s.Participating = ss.Participating

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", s, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Store) ToSnapshot() es.Snapshot {
	return StoreV1{
		Name:          s.Name,
		Location:      s.Location,
		Participating: s.Participating,
	}
}
