package jetstream

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	"eda-in-golang/internal/am"
)

const maxRetries = 5

type Stream struct {
	streamName string
	js         nats.JetStreamContext
	mu         sync.Mutex
	logger     zerolog.Logger
}

var _ am.RawMessageStream = (*Stream)(nil)

func NewStream(streamName string, js nats.JetStreamContext, logger zerolog.Logger) *Stream {
	return &Stream{
		streamName: streamName,
		js:         js,
		logger:     logger,
	}
}

func (s *Stream) Publish(ctx context.Context, topicName string, rawMsg am.RawMessage) (err error) {
	var data []byte

	data, err = proto.Marshal(&StreamMessage{
		Id:   rawMsg.ID(),
		Name: rawMsg.MessageName(),
		Data: rawMsg.Data(),
	})
	if err != nil {
		return
	}

	var p nats.PubAckFuture
	p, err = s.js.PublishMsgAsync(&nats.Msg{
		Subject: topicName,
		Data:    data,
	}, nats.MsgId(rawMsg.ID()))
	if err != nil {
		return
	}

	// retry a handful of times to publish the messages
	go func(future nats.PubAckFuture, tries int) {
		var err error

		for {
			select {
			case <-future.Ok(): // publish acknowledged
				return
			case <-future.Err(): // error ignored; try again
				// TODO add some variable delay between tries
				tries = tries - 1
				if tries <= 0 {
					s.logger.Error().Msgf("unable to publish message after %d tries", maxRetries)
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					// TODO do more than give up
					s.logger.Error().Err(err).Msg("failed to publish a message")
					return
				}
			}
		}
	}(p, maxRetries)

	return
}

func (s *Stream) Subscribe(topicName string, handler am.RawMessageHandler, options ...am.SubscriberOption) error {
	var err error

	s.mu.Lock()
	defer s.mu.Unlock()

	subCfg := am.NewSubscriberConfig(options)

	opts := []nats.SubOpt{
		nats.MaxDeliver(subCfg.MaxRedeliver()),
	}
	cfg := &nats.ConsumerConfig{
		MaxDeliver:     subCfg.MaxRedeliver(),
		DeliverSubject: topicName,
		FilterSubject:  topicName,
	}
	if groupName := subCfg.GroupName(); groupName != "" {
		cfg.DeliverSubject = groupName
		cfg.DeliverGroup = groupName
		cfg.Durable = groupName

		opts = append(opts, nats.Bind(s.streamName, groupName), nats.Durable(groupName))
	}

	if ackType := subCfg.AckType(); ackType != am.AckTypeAuto {
		ackWait := subCfg.AckWait()

		cfg.AckPolicy = nats.AckExplicitPolicy
		cfg.AckWait = ackWait

		opts = append(opts, nats.AckExplicit(), nats.AckWait(ackWait))
	} else {
		cfg.AckPolicy = nats.AckNonePolicy
		opts = append(opts, nats.AckNone())
	}

	_, err = s.js.AddConsumer(s.streamName, cfg)
	if err != nil {
		return err
	}

	if groupName := subCfg.GroupName(); groupName == "" {
		_, err = s.js.Subscribe(topicName, s.handleMsg(subCfg, handler), opts...)
	} else {
		_, err = s.js.QueueSubscribe(topicName, groupName, s.handleMsg(subCfg, handler), opts...)
	}

	return nil
}

func (s *Stream) handleMsg(cfg am.SubscriberConfig, handler am.RawMessageHandler) func(*nats.Msg) {
	return func(natsMsg *nats.Msg) {
		var err error

		m := &StreamMessage{}
		err = proto.Unmarshal(natsMsg.Data, m)
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to unmarshal the *nats.Msg")
			return
		}

		msg := &rawMessage{
			id:       m.GetId(),
			name:     m.GetName(),
			data:     m.GetData(),
			acked:    false,
			ackFn:    func() error { return natsMsg.Ack() },
			nackFn:   func() error { return natsMsg.Nak() },
			extendFn: func() error { return natsMsg.InProgress() },
			killFn:   func() error { return natsMsg.Term() },
		}

		wCtx, cancel := context.WithTimeout(context.Background(), cfg.AckWait())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err = msg.Ack()
			if err != nil {
				s.logger.Warn().Err(err).Msg("failed to auto-Ack a message")
			}
		}

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := msg.Ack(); ackErr != nil {
					s.logger.Warn().Err(err).Msg("failed to Ack a message")
				}
				return
			}
			s.logger.Error().Err(err).Msg("error while handling message")
			if nakErr := msg.NAck(); nakErr != nil {
				s.logger.Warn().Err(err).Msg("failed to Nack a message")
			}
		case <-wCtx.Done():
			// TODO logging?
			return
		}
	}
}
