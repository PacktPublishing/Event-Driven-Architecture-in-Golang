package ddd

import (
	"context"
)

type EventHandler func(ctx context.Context, event Event) error

type Event interface {
	EventName() string
}
