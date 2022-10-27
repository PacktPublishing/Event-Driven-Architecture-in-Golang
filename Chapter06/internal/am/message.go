package am

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type (
	Message interface {
		ddd.IDer
		MessageName() string
		Ack() error
		NAck() error
		Extend() error
		Kill() error
	}

	MessageHandler[O Message] interface {
		HandleMessage(ctx context.Context, msg O) error
	}

	MessageHandlerFunc[O Message] func(ctx context.Context, msg O) error

	MessagePublisher[I any] interface {
		Publish(ctx context.Context, topicName string, v I) error
	}

	MessageSubscriber[O Message] interface {
		Subscribe(topicName string, handler MessageHandler[O], options ...SubscriberOption) error
	}

	MessageStream[I any, O Message] interface {
		MessagePublisher[I]
		MessageSubscriber[O]
	}
)

func (f MessageHandlerFunc[O]) HandleMessage(ctx context.Context, msg O) error {
	return f(ctx, msg)
}
