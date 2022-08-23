package am

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

type (
	EventMessage interface {
		Message
		ddd.Event
	}

	EventPublisher  = MessagePublisher[ddd.Event]
	EventSubscriber = MessageSubscriber[EventMessage]
	EventStream     = MessageStream[ddd.Event, EventMessage]

	eventStream struct {
		reg    registry.Registry
		stream MessageStream[RawMessage, RawMessage]
	}

	eventMessage struct {
		id         string
		name       string
		payload    ddd.EventPayload
		metadata   ddd.Metadata
		occurredAt time.Time
		msg        RawMessage
	}
)

var _ EventMessage = (*eventMessage)(nil)

var _ EventStream = (*eventStream)(nil)

func NewEventStream(reg registry.Registry, stream MessageStream[RawMessage, RawMessage]) EventStream {
	return &eventStream{
		reg:    reg,
		stream: stream,
	}
}

func (s eventStream) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	metadata, err := structpb.NewStruct(event.Metadata())
	if err != nil {
		return err
	}

	payload, err := s.reg.Serialize(
		event.EventName(), event.Payload(),
	)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&EventMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(event.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:   event.ID(),
		name: event.EventName(),
		data: data,
	})
}

func (s eventStream) Subscribe(topicName string, handler MessageHandler[EventMessage], options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[RawMessage](func(ctx context.Context, msg RawMessage) error {
		var eventData EventMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &eventData)
		if err != nil {
			return err
		}

		eventName := msg.MessageName()

		payload, err := s.reg.Deserialize(eventName, eventData.GetPayload())
		if err != nil {
			return err
		}

		eventMsg := eventMessage{
			id:         msg.ID(),
			name:       eventName,
			payload:    payload,
			metadata:   eventData.GetMetadata().AsMap(),
			occurredAt: eventData.GetOccurredAt().AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, eventMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (e eventMessage) ID() string                { return e.id }
func (e eventMessage) EventName() string         { return e.name }
func (e eventMessage) Payload() ddd.EventPayload { return e.payload }
func (e eventMessage) Metadata() ddd.Metadata    { return e.metadata }
func (e eventMessage) OccurredAt() time.Time     { return e.occurredAt }
func (e eventMessage) MessageName() string       { return e.msg.MessageName() }
func (e eventMessage) Ack() error                { return e.msg.Ack() }
func (e eventMessage) NAck() error               { return e.msg.NAck() }
func (e eventMessage) Extend() error             { return e.msg.Extend() }
func (e eventMessage) Kill() error               { return e.msg.Kill() }
