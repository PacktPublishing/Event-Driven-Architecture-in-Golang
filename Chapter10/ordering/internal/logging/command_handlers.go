package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/ddd"
)

type CommandHandlers[T ddd.Command] struct {
	ddd.CommandHandler[T]
	label  string
	logger zerolog.Logger
}

var _ ddd.CommandHandler[ddd.Command] = (*CommandHandlers[ddd.Command])(nil)

func LogCommandHandlerAccess[T ddd.Command](handlers ddd.CommandHandler[T], label string, logger zerolog.Logger) ddd.CommandHandler[T] {
	return CommandHandlers[T]{
		CommandHandler: handlers,
		label:          label,
		logger:         logger,
	}
}

func (h CommandHandlers[T]) HandleCommand(ctx context.Context, command T) (reply ddd.Reply, err error) {
	h.logger.Info().Msgf("--> Ordering.%s.On(%s)", h.label, command.CommandName())
	defer func() { h.logger.Info().Err(err).Msgf("<-- Ordering.%s.On(%s)", h.label, command.CommandName()) }()
	return h.CommandHandler.HandleCommand(ctx, command)
}
