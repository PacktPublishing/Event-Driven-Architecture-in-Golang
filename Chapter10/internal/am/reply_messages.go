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
	ReplyMessage interface {
		Message
		ddd.Reply
	}

	IncomingReplyMessage interface {
		IncomingMessage
		ddd.Reply
	}

	ReplyPublisher  = MessagePublisher[ddd.Reply]
	ReplySubscriber = MessageSubscriber[IncomingReplyMessage]
	ReplyStream     = MessageStream[ddd.Reply, IncomingReplyMessage]

	replyStream struct {
		reg    registry.Registry
		stream RawMessageStream
	}

	replyMessage struct {
		id         string
		name       string
		payload    ddd.ReplyPayload
		metadata   ddd.Metadata
		occurredAt time.Time
		msg        IncomingMessage
	}
)

var _ ReplyMessage = (*replyMessage)(nil)

var _ ReplyStream = (*replyStream)(nil)

func NewReplyStream(reg registry.Registry, stream RawMessageStream) ReplyStream {
	return &replyStream{
		reg:    reg,
		stream: stream,
	}
}

func (s replyStream) Publish(ctx context.Context, topicName string, reply ddd.Reply) error {
	metadata, err := structpb.NewStruct(reply.Metadata())
	if err != nil {
		return err
	}

	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = s.reg.Serialize(reply.ReplyName(), reply.Payload())
		if err != nil {
			return err
		}
	}

	data, err := proto.Marshal(&ReplyMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(reply.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:      reply.ID(),
		name:    reply.ReplyName(),
		subject: topicName,
		data:    data,
	})
}

func (s replyStream) Subscribe(topicName string, handler MessageHandler[IncomingReplyMessage], options ...SubscriberOption) (Subscription, error) {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[IncomingRawMessage](func(ctx context.Context, msg IncomingRawMessage) error {
		var replyData ReplyMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &replyData)
		if err != nil {
			return err
		}

		replyName := msg.MessageName()

		var payload any

		if replyName != SuccessReply && replyName != FailureReply {
			payload, err = s.reg.Deserialize(replyName, replyData.GetPayload())
			if err != nil {
				return err
			}
		}

		replyMsg := replyMessage{
			id:         msg.ID(),
			name:       replyName,
			payload:    payload,
			metadata:   replyData.GetMetadata().AsMap(),
			occurredAt: replyData.GetOccurredAt().AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, replyMsg)
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (s replyStream) Unsubscribe() error {
	return s.stream.Unsubscribe()
}

func (r replyMessage) ID() string                { return r.id }
func (r replyMessage) ReplyName() string         { return r.name }
func (r replyMessage) Payload() ddd.ReplyPayload { return r.payload }
func (r replyMessage) Metadata() ddd.Metadata    { return r.metadata }
func (r replyMessage) OccurredAt() time.Time     { return r.occurredAt }
func (r replyMessage) Subject() string           { return r.msg.Subject() }
func (r replyMessage) MessageName() string       { return r.msg.MessageName() }
func (r replyMessage) Ack() error                { return r.msg.Ack() }
func (r replyMessage) NAck() error               { return r.msg.NAck() }
func (r replyMessage) Extend() error             { return r.msg.Extend() }
func (r replyMessage) Kill() error               { return r.msg.Kill() }

type replyMsgHandler struct {
	reg     registry.Registry
	handler ddd.ReplyHandler[ddd.Reply]
}

func NewReplyMessageHandler(reg registry.Registry, handler ddd.ReplyHandler[ddd.Reply]) RawMessageHandler {
	return replyMsgHandler{
		reg:     reg,
		handler: handler,
	}
}

func (h replyMsgHandler) HandleMessage(ctx context.Context, msg IncomingRawMessage) error {
	var replyData ReplyMessageData

	err := proto.Unmarshal(msg.Data(), &replyData)
	if err != nil {
		return err
	}

	replyName := msg.MessageName()

	var payload any

	if replyName != SuccessReply && replyName != FailureReply {
		payload, err = h.reg.Deserialize(replyName, replyData.GetPayload())
		if err != nil {
			return err
		}
	}

	replyMsg := replyMessage{
		id:         msg.ID(),
		name:       replyName,
		payload:    payload,
		metadata:   replyData.GetMetadata().AsMap(),
		occurredAt: replyData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	return h.handler.HandleReply(ctx, replyMsg)
}
