package am

import (
	"context"

	"eda-in-golang/internal/ddd"
)

const (
	CommandHdrPrefix       = "COMMAND_"
	CommandNameHdr         = CommandHdrPrefix + "NAME"
	CommandReplyChannelHdr = CommandHdrPrefix + "REPLY_CHANNEL"
)

type (
	CommandMessageHandler interface {
		HandleMessage(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error)
	}

	CommandMessageHandlerFunc func(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error)

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

func (f CommandMessageHandlerFunc) HandleMessage(ctx context.Context, cmd IncomingCommandMessage) (ddd.Reply, error) {
	return f(ctx, cmd)
}
