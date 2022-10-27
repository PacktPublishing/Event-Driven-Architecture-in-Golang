package am

import (
	"context"
)

type (
	RawMessageStream           = MessageStream[RawMessage, IncomingRawMessage]
	RawMessagePublisher        = MessagePublisher[RawMessage]
	RawMessageSubscriber       = MessageSubscriber[IncomingRawMessage]
	RawMessageStreamMiddleware = func(stream RawMessageStream) RawMessageStream

	RawMessageHandler           = MessageHandler[IncomingRawMessage]
	RawMessageHandlerFunc       func(ctx context.Context, msg IncomingRawMessage) error
	RawMessageHandlerMiddleware = func(handler RawMessageHandler) RawMessageHandler

	RawMessage interface {
		Message
		Data() []byte
	}

	IncomingRawMessage interface {
		IncomingMessage
		Data() []byte
	}

	rawMessage struct {
		id      string
		name    string
		subject string
		data    []byte
	}
)

var _ RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string          { return m.id }
func (m rawMessage) Subject() string     { return m.subject }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte        { return m.data }

func (f RawMessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingRawMessage) error {
	return f(ctx, cmd)
}

func RawMessageStreamWithMiddleware(stream RawMessageStream, mws ...RawMessageStreamMiddleware) RawMessageStream {
	s := stream
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}
	return s
}

func RawMessageHandlerWithMiddleware(handler RawMessageHandler, mws ...RawMessageHandlerMiddleware) RawMessageHandler {
	h := handler
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
