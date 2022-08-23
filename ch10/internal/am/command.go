package am

import (
	"context"
	"strings"

	"google.golang.org/protobuf/proto"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

const (
	CommandHdrPrefix       = "COMMAND_"
	CommandNameHdr         = CommandHdrPrefix + "NAME"
	CommandReplyChannelHdr = CommandHdrPrefix + "REPLY_CHANNEL"
)

type (
	CommandMessageHandler     = MessageHandler[IncomingCommandMessage]
	CommandMessageHandlerFunc func(ctx context.Context, msg IncomingCommandMessage) error

	Command interface {
		ddd.Command
		Destination() string
	}
	command struct {
		ddd.Command
		destination string
	}
)

func NewCommand(name, destination string, payload ddd.CommandPayload, options ...ddd.CommandOption) Command {
	return command{
		Command:     ddd.NewCommand(name, payload, options...),
		destination: destination,
	}
}

func (c command) Destination() string {
	return c.destination
}

func (f CommandMessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingCommandMessage) error {
	return f(ctx, cmd)
}

type commandMsgHandler struct {
	reg       registry.Registry
	publisher ReplyPublisher
	handler   ddd.CommandHandler[ddd.Command]
}

func NewCommandMessageHandler(reg registry.Registry, publisher ReplyPublisher, handler ddd.CommandHandler[ddd.Command]) RawMessageHandler {
	return commandMsgHandler{
		reg:       reg,
		publisher: publisher,
		handler:   handler,
	}
}

func (h commandMsgHandler) HandleMessage(ctx context.Context, msg IncomingRawMessage) error {
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
		metadata:   commandData.GetMetadata().AsMap(),
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
