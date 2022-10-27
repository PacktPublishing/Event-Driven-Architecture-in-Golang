package am

import (
	"context"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

const (
	CommandHdrPrefix       = "COMMAND_"
	CommandNameHdr         = CommandHdrPrefix + "NAME"
	CommandReplyChannelHdr = CommandHdrPrefix + "REPLY_CHANNEL"
)

type (
	CommandMessage interface {
		MessageBase
		ddd.Command
	}

	IncomingCommandMessage interface {
		IncomingMessageBase
		ddd.Command
	}

	CommandPublisher interface {
		Publish(ctx context.Context, topicName string, cmd ddd.Command) error
	}

	commandPublisher struct {
		reg       registry.Registry
		publisher MessagePublisher
	}

	commandMessage struct {
		id         string
		name       string
		payload    ddd.CommandPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}

	commandMsgHandler struct {
		reg       registry.Registry
		publisher ReplyPublisher
		handler   ddd.CommandHandler[ddd.Command]
	}
)

var _ CommandMessage = (*commandMessage)(nil)

var _ CommandPublisher = (*commandPublisher)(nil)

func NewCommandPublisher(reg registry.Registry, msgPublisher MessagePublisher, mws ...MessagePublisherMiddleware) CommandPublisher {
	return commandPublisher{
		reg:       reg,
		publisher: MessagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

func (s commandPublisher) Publish(ctx context.Context, topicName string, command ddd.Command) error {
	payload, err := s.reg.Serialize(command.CommandName(), command.Payload())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(command.OccurredAt()),
	})
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, topicName, message{
		id:       command.ID(),
		name:     command.CommandName(),
		subject:  topicName,
		data:     data,
		metadata: command.Metadata(),
		sentAt:   time.Now(),
	})
}

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Payload() ddd.CommandPayload { return c.payload }
func (c commandMessage) Metadata() ddd.Metadata      { return c.msg.Metadata() }
func (c commandMessage) OccurredAt() time.Time       { return c.occurredAt }
func (c commandMessage) Subject() string             { return c.msg.Subject() }
func (c commandMessage) MessageName() string         { return c.msg.MessageName() }
func (c commandMessage) SentAt() time.Time           { return c.msg.SentAt() }
func (c commandMessage) ReceivedAt() time.Time       { return c.msg.ReceivedAt() }
func (c commandMessage) Ack() error                  { return c.msg.Ack() }
func (c commandMessage) NAck() error                 { return c.msg.NAck() }
func (c commandMessage) Extend() error               { return c.msg.Extend() }
func (c commandMessage) Kill() error                 { return c.msg.Kill() }

func NewCommandHandler(reg registry.Registry, publisher ReplyPublisher, handler ddd.CommandHandler[ddd.Command], mws ...MessageHandlerMiddleware) MessageHandler {
	return MessageHandlerWithMiddleware(commandMsgHandler{
		reg:       reg,
		publisher: publisher,
		handler:   handler,
	}, mws...)
}

func (h commandMsgHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	var commandData CommandMessageData

	err := proto.Unmarshal(msg.Data(), &commandData)
	if err != nil {
		return err
	}

	commandName := msg.MessageName()

	payload, err := h.reg.Deserialize(commandName, commandData.GetPayload())
	if err != nil {
		return err
	}

	commandMsg := commandMessage{
		id:         msg.ID(),
		name:       commandName,
		payload:    payload,
		occurredAt: commandData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	destination := commandMsg.Metadata().Get(CommandReplyChannelHdr).(string)

	reply, err := h.handler.HandleCommand(ctx, commandMsg)
	if err != nil {
		return h.publishReply(ctx, destination, h.failure(reply, commandMsg))
	}

	return h.publishReply(ctx, destination, h.success(reply, commandMsg))
}

func (h commandMsgHandler) publishReply(ctx context.Context, destination string, reply ddd.Reply) error {
	return h.publisher.Publish(ctx, destination, reply)
}

func (h commandMsgHandler) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeFailure)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMsgHandler) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeSuccess)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMsgHandler) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	for key, value := range cmd.Metadata() {
		if key == CommandNameHdr {
			continue
		}

		if strings.HasPrefix(key, CommandHdrPrefix) {
			hdr := ReplyHdrPrefix + key[len(CommandHdrPrefix):]
			reply.Metadata().Set(hdr, value)
		}
	}

	return reply
}
