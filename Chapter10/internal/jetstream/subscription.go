package jetstream

import (
	"github.com/nats-io/nats.go"
)

type subscription struct {
	s *nats.Subscription
}

func (s subscription) Unsubscribe() error {
	if !s.s.IsValid() {
		return nil
	}

	return s.s.Drain()
}
