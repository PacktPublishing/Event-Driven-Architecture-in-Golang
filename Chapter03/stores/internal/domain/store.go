package domain

import (
	"github.com/stackus/errors"
)

var (
	ErrStoreNameIsBlank               = errors.Wrap(errors.ErrBadRequest, "the store name cannot be blank")
	ErrStoreLocationIsBlank           = errors.Wrap(errors.ErrBadRequest, "the store location cannot be blank")
	ErrStoreIsAlreadyParticipating    = errors.Wrap(errors.ErrBadRequest, "the store is already participating")
	ErrStoreIsAlreadyNotParticipating = errors.Wrap(errors.ErrBadRequest, "the store is already not participating")
)

type Store struct {
	ID            string
	Name          string
	Location      string
	Participating bool
}

func CreateStore(id, name, location string) (store *Store, err error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	if location == "" {
		return nil, ErrStoreLocationIsBlank
	}

	store = &Store{
		ID:       id,
		Name:     name,
		Location: location,
	}

	return
}

func (s *Store) EnableParticipation() (err error) {
	if s.Participating {
		return ErrStoreIsAlreadyParticipating
	}

	s.Participating = true

	return
}

func (s *Store) DisableParticipation() (err error) {
	if !s.Participating {
		return ErrStoreIsAlreadyNotParticipating
	}

	s.Participating = false

	return
}
