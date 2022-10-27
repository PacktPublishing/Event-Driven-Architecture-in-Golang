package jetstream

import (
	"eda-in-golang/internal/am"
)

type rawMessage struct {
	id       string
	name     string
	subject  string
	data     []byte
	acked    bool
	ackFn    func() error
	nackFn   func() error
	extendFn func() error
	killFn   func() error
}

var _ am.RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string          { return m.id }
func (m rawMessage) Subject() string     { return m.subject }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte        { return m.data }

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

func (m rawMessage) Extend() error {
	return m.extendFn()
}

func (m *rawMessage) Kill() error {
	if m.acked {
		return nil
	}

	m.acked = true
	return m.killFn()
}
