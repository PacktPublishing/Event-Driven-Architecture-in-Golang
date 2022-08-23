package domain

import (
	"context"
)

type ParticipatingStoreRepository interface {
	FindAll(ctx context.Context) ([]*Store, error)
}
