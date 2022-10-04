package tm

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/internal/am"
)

type OutboxStore interface {
	Save(ctx context.Context, msg am.Message) error
	FindUnpublished(ctx context.Context, limit int) ([]am.Message, error)
	MarkPublished(ctx context.Context, ids ...string) error
}

func OutboxPublisher(store OutboxStore) am.MessagePublisherMiddleware {
	return func(next am.MessagePublisher) am.MessagePublisher {
		return am.MessagePublisherFunc(func(ctx context.Context, topicName string, msg am.Message) error {
			err := store.Save(ctx, msg)
			var errDupe ErrDuplicateMessage
			if errors.As(err, &errDupe) {
				return nil
			}
			return err
		})
	}
}
