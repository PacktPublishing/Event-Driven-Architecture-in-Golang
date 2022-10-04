package sec

import (
	"context"

	"eda-in-golang/internal/ddd"
)

type (
	StepActionFunc[T any]       func(ctx context.Context, data T) (string, ddd.Command, error)
	StepReplyHandlerFunc[T any] func(ctx context.Context, data T, reply ddd.Reply) error

	SagaStep[T any] interface {
		Action(fn StepActionFunc[T]) SagaStep[T]
		Compensation(fn StepActionFunc[T]) SagaStep[T]
		OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]
		OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]
		isInvocable(compensating bool) bool
		execute(ctx context.Context, sagaCtx *SagaContext[T]) stepResult[T]
		handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error
	}

	sagaStep[T any] struct {
		actions  map[bool]StepActionFunc[T]
		handlers map[bool]map[string]StepReplyHandlerFunc[T]
	}

	stepResult[T any] struct {
		ctx         *SagaContext[T]
		destination string
		cmd         ddd.Command
		err         error
	}
)

var _ SagaStep[any] = (*sagaStep[any])(nil)

func (s *sagaStep[T]) Action(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[notCompensating] = fn
	return s
}

func (s *sagaStep[T]) Compensation(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[isCompensating] = fn
	return s
}

func (s *sagaStep[T]) OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[notCompensating][replyName] = fn
	return s
}

func (s *sagaStep[T]) OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s sagaStep[T]) isInvocable(compensating bool) bool {
	return s.actions[compensating] != nil
}

func (s sagaStep[T]) execute(ctx context.Context, sagaCtx *SagaContext[T]) stepResult[T] {
	if action := s.actions[sagaCtx.Compensating]; action != nil {
		destination, cmd, err := action(ctx, sagaCtx.Data)
		return stepResult[T]{
			ctx:         sagaCtx,
			destination: destination,
			cmd:         cmd,
			err:         err,
		}
	}

	return stepResult[T]{ctx: sagaCtx}
}

func (s sagaStep[T]) handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error {
	if handler := s.handlers[sagaCtx.Compensating][reply.ReplyName()]; handler != nil {
		return handler(ctx, sagaCtx.Data, reply)
	}
	return nil
}

type StepOption[T any] func(step *sagaStep[T])

func WithAction[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[notCompensating] = fn
	}
}

func WithCompensation[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[isCompensating] = fn
	}
}

func OnActionReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.handlers[notCompensating][replyName] = fn
	}
}

func OnCompensationReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.handlers[isCompensating][replyName] = fn
	}
}
