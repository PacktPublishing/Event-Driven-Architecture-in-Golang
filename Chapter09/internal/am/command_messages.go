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
	CommandMessage interface {
		Message
		ddd.Command
	}

	IncomingCommandMessage interface {
		IncomingMessage
		ddd.Command
	}

	CommandPublisher  = MessagePublisher[ddd.Command]
	CommandSubscriber = MessageSubscriber[IncomingCommandMessage]
	CommandStream     interface {
		MessagePublisher[ddd.Command]
		MessageSubscriber[IncomingCommandMessage]
	}

	commandStream struct {
		reg    registry.Registry
		stream RawMessageStream
	}

	commandMessage struct {
		id         string
		name       string
		payload    ddd.CommandPayload
		metadata   ddd.Metadata
		occurredAt time.Time
		msg        IncomingMessage
	}
)

var _ CommandMessage = (*commandMessage)(nil)

var _ CommandStream = (*commandStream)(nil)

func NewCommandStream(reg registry.Registry, stream RawMessageStream) CommandStream {
	return commandStream{
		reg:    reg,
		stream: stream,
	}
}

func (s commandStream) Publish(ctx context.Context, topicName string, command ddd.Command) error {
	metadata, err := structpb.NewStruct(command.Metadata())
	if err != nil {
		return err
	}

	payload, err := s.reg.Serialize(
		command.CommandName(), command.Payload(),
	)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(command.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:      command.ID(),
		name:    command.CommandName(),
		subject: topicName,
		data:    data,
	})
}

func (s commandStream) Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[IncomingRawMessage](func(ctx context.Context, msg IncomingRawMessage) error {
		var commandData CommandMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &commandData)
		if err != nil {
			return err
		}

		commandName := msg.MessageName()

		payload, err := s.reg.Deserialize(commandName, commandData.GetPayload())
		if err != nil {
			return err
		}

		commandMsg := commandMessage{
			id:         msg.ID(),
			name:       commandName,
			payload:    payload,
			metadata:   commandData.GetMetadata().AsMap(),
			occurredAt: commandData.GetOccurredAt().AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, commandMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Payload() ddd.CommandPayload { return c.payload }
func (c commandMessage) Metadata() ddd.Metadata      { return c.metadata }
func (c commandMessage) OccurredAt() time.Time       { return c.occurredAt }
func (c commandMessage) Subject() string             { return c.msg.Subject() }
func (c commandMessage) MessageName() string         { return c.msg.MessageName() }
func (c commandMessage) Ack() error                  { return c.msg.Ack() }
func (c commandMessage) NAck() error                 { return c.msg.NAck() }
func (c commandMessage) Extend() error               { return c.msg.Extend() }
func (c commandMessage) Kill() error                 { return c.msg.Kill() }
