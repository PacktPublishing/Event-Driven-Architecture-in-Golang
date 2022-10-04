package amprom

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"eda-in-golang/internal/am"
)

func SentMessagesCounter(serviceName string) am.MessagePublisherMiddleware {
	counter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Name:      "sent_messages_count",
		Help:      fmt.Sprintf("The total number of messages sent by %s", serviceName),
	}, []string{"message"})

	return func(next am.MessagePublisher) am.MessagePublisher {
		return am.MessagePublisherFunc(func(ctx context.Context, topicName string, msg am.Message) (err error) {
			defer func() {
				counter.WithLabelValues("all").Inc()
				counter.WithLabelValues(msg.MessageName()).Inc()
			}()
			return next.Publish(ctx, topicName, msg)
		})
	}
}
