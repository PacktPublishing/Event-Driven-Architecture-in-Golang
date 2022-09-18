package am

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type (
	Message interface {
		ddd.IDer
		Subject() string
		MessageName() string
	}

	IncomingMessage interface {
		Message
		Ack() error
		NAck() error
		Extend() error
		Kill() error
	}

	MessageHandler[I IncomingMessage] interface {
		HandleMessage(ctx context.Context, msg I) error
	}

	MessageHandlerFunc[I IncomingMessage] func(ctx context.Context, msg I) error

	MessagePublisher[O any] interface {
		Publish(ctx context.Context, topicName string, v O) error
	}

	MessageSubscriber[I IncomingMessage] interface {
		Subscribe(topicName string, handler MessageHandler[I], options ...SubscriberOption) (Subscription, error)
		Unsubscribe() error
	}

	MessageStream[O any, I IncomingMessage] interface {
		MessagePublisher[O]
		MessageSubscriber[I]
	}
)

func (f MessageHandlerFunc[I]) HandleMessage(ctx context.Context, msg I) error {
	return f(ctx, msg)
}
