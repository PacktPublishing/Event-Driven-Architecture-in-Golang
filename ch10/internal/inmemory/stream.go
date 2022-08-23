package inmemory

import (
	"context"

	"eda-in-golang/internal/am"
)

type stream struct {
	subscriptions map[string][]am.RawMessageHandler
}

var _ am.RawMessageStream = (*stream)(nil)

func NewStream() stream {
	return stream{
		subscriptions: make(map[string][]am.RawMessageHandler),
	}
}

func (t stream) Publish(ctx context.Context, topicName string, v am.RawMessage) error {
	for _, handler := range t.subscriptions[topicName] {
		err := handler.HandleMessage(ctx, &rawMessage{v})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t stream) Subscribe(topicName string, handler am.RawMessageHandler, options ...am.SubscriberOption) error {
	cfg := am.NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := am.RawMessageHandlerFunc(func(ctx context.Context, msg am.IncomingRawMessage) error {
		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		return handler.HandleMessage(ctx, msg)
	})

	t.subscriptions[topicName] = append(t.subscriptions[topicName], fn)

	return nil
}
