package amotel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"eda-in-golang/internal/am"
)

func OtelMessageContextInjector() am.MessagePublisherMiddleware {
	return func(next am.MessagePublisher) am.MessagePublisher {
		return am.MessagePublisherFunc(func(ctx context.Context, topicName string, msg am.Message) error {
			var span trace.Span
			ctx, span = tracer.Start(ctx,
				msg.MessageName(),
				trace.WithSpanKind(trace.SpanKindProducer),
				trace.WithAttributes(
					attribute.String("MessageID", msg.ID()),
					attribute.String("Subject", msg.Subject()),
				),
			)
			defer span.End()
			propagator.Inject(ctx, MetadataCarrier(msg.Metadata()))

			err := next.Publish(ctx, topicName, msg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			return err
		})
	}
}
