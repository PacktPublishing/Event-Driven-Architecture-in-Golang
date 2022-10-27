package am

import (
	"context"

	"github.com/stackus/errors"
)

type fakeMessage[O any] struct {
	subject string
	payload O
}

type FakeMessagePublisher[O any] struct {
	messages []fakeMessage[O]
}

var _ MessagePublisher[any] = (*FakeMessagePublisher[any])(nil)

func NewFakeMessagePublisher[O any]() *FakeMessagePublisher[O] {
	return &FakeMessagePublisher[O]{
		messages: []fakeMessage[O]{},
	}
}

func (p *FakeMessagePublisher[O]) Publish(ctx context.Context, topicName string, v O) error {
	p.messages = append(p.messages, fakeMessage[O]{topicName, v})
	return nil
}

func (p *FakeMessagePublisher[O]) Reset() {
	p.messages = []fakeMessage[O]{}
}

func (p *FakeMessagePublisher[O]) Last() (string, O, error) {
	var v O
	if len(p.messages) == 0 {
		return "", v, errors.ErrNotFound.Msg("no messages have been published")
	}

	last := p.messages[len(p.messages)-1]
	return last.subject, last.payload, nil
}
