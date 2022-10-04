package am

import (
	"context"
	"time"

	"eda-in-golang/internal/ddd"
)

type (
	MessageBase interface {
		ddd.IDer
		Subject() string
		MessageName() string
		Metadata() ddd.Metadata
		SentAt() time.Time
	}

	IncomingMessageBase interface {
		MessageBase
		ReceivedAt() time.Time
		Ack() error
		NAck() error
		Extend() error
		Kill() error
	}

	MessageHandler interface {
		HandleMessage(ctx context.Context, msg IncomingMessage) error
	}
	MessageHandlerFunc func(ctx context.Context, msg IncomingMessage) error

	MessageSubscriber interface {
		Subscribe(topicName string, handler MessageHandler, options ...SubscriberOption) (Subscription, error)
		Unsubscribe() error
	}

	MessagePublisher interface {
		Publish(ctx context.Context, topicName string, msg Message) error
	}
	MessageStream interface {
		MessageSubscriber
		MessagePublisher
	}

	MessagePublisherFunc func(ctx context.Context, topicName string, msg Message) error

	MessageStreamMiddleware    = func(next MessageStream) MessageStream
	MessagePublisherMiddleware = func(next MessagePublisher) MessagePublisher
	MessageHandlerMiddleware   = func(next MessageHandler) MessageHandler

	Message interface {
		MessageBase
		Data() []byte
	}

	IncomingMessage interface {
		IncomingMessageBase
		Data() []byte
	}

	message struct {
		id       string
		name     string
		subject  string
		data     []byte
		metadata ddd.Metadata
		sentAt   time.Time
	}

	messagePublisher struct {
		publisher MessagePublisher
	}

	messageSubscriber struct {
		subscriber MessageSubscriber
		mws        []MessageHandlerMiddleware
	}
)

var _ Message = (*message)(nil)

func (m message) ID() string             { return m.id }
func (m message) Subject() string        { return m.subject }
func (m message) MessageName() string    { return m.name }
func (m message) Data() []byte           { return m.data }
func (m message) Metadata() ddd.Metadata { return m.metadata }
func (m message) SentAt() time.Time      { return m.sentAt }

func (f MessagePublisherFunc) Publish(ctx context.Context, topicName string, msg Message) error {
	return f(ctx, topicName, msg)
}

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingMessage) error {
	return f(ctx, cmd)
}

func NewMessagePublisher(publisher MessagePublisher, mws ...MessagePublisherMiddleware) MessagePublisher {
	return messagePublisher{
		publisher: MessagePublisherWithMiddleware(publisher, mws...),
	}
}

func (p messagePublisher) Publish(ctx context.Context, topicName string, msg Message) error {
	return p.publisher.Publish(ctx, topicName, msg)
}

func NewMessageSubscriber(subscriber MessageSubscriber, mws ...MessageHandlerMiddleware) MessageSubscriber {
	return messageSubscriber{
		subscriber: subscriber,
		mws:        mws,
	}
}

func (s messageSubscriber) Subscribe(topicName string, handler MessageHandler, options ...SubscriberOption) (Subscription, error) {
	return s.subscriber.Subscribe(topicName, MessageHandlerWithMiddleware(handler, s.mws...), options...)
}

func (s messageSubscriber) Unsubscribe() error {
	return s.subscriber.Unsubscribe()
}
