package jetstream

import (
	"time"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type rawMessage struct {
	id         string
	name       string
	subject    string
	data       []byte
	metadata   ddd.Metadata
	sentAt     time.Time
	receivedAt time.Time
	acked      bool
	ackFn      func() error
	nackFn     func() error
	extendFn   func() error
	killFn     func() error
}

var _ am.Message = (*rawMessage)(nil)

func (m *rawMessage) ID() string             { return m.id }
func (m *rawMessage) Subject() string        { return m.subject }
func (m *rawMessage) MessageName() string    { return m.name }
func (m *rawMessage) Data() []byte           { return m.data }
func (m *rawMessage) Metadata() ddd.Metadata { return m.metadata }
func (m *rawMessage) SentAt() time.Time      { return m.sentAt }
func (m *rawMessage) ReceivedAt() time.Time  { return m.receivedAt }

func (m *rawMessage) Ack() error {
	if m.acked {
		return nil
	}
	m.acked = true
	return m.ackFn()
}

func (m *rawMessage) NAck() error {
	if m.acked {
		return nil
	}
	m.acked = true
	return m.nackFn()
}

func (m *rawMessage) Extend() error {
	return m.extendFn()
}

func (m *rawMessage) Kill() error {
	if m.acked {
		return nil
	}

	m.acked = true
	return m.killFn()
}
