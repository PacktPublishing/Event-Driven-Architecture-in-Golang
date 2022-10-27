package tm

import (
	"context"
	"time"

	"eda-in-golang/internal/am"
)

const messageLimit = 50
const pollingInterval = 500 * time.Millisecond

type OutboxProcessor interface {
	Start(ctx context.Context) error
}

type outboxProcessor struct {
	publisher am.RawMessagePublisher
	store     OutboxStore
}

func NewOutboxProcessor(publisher am.RawMessagePublisher, store OutboxStore) OutboxProcessor {
	return outboxProcessor{
		publisher: publisher,
		store:     store,
	}
}

func (p outboxProcessor) Start(ctx context.Context) error {
	errC := make(chan error)

	go func() {
		errC <- p.processMessages(ctx)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errC:
		return err
	}
}

func (p outboxProcessor) processMessages(ctx context.Context) error {
	timer := time.NewTimer(0)
	for {
		msgs, err := p.store.FindUnpublished(ctx, messageLimit)
		if err != nil {
			return err
		}

		if len(msgs) > 0 {
			ids := make([]string, len(msgs))
			for i, msg := range msgs {
				ids[i] = msg.ID()
				err = p.publisher.Publish(ctx, msg.Subject(), msg)
				if err != nil {
					return err
				}
			}
			err = p.store.MarkPublished(ctx, ids...)
			if err != nil {
				return err
			}

			// poll again immediately
			continue
		}

		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}

		// wait a short time before polling again
		timer.Reset(pollingInterval)

		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
		}
	}
}
