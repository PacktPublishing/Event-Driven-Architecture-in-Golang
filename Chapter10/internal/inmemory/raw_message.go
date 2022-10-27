package inmemory

import (
	"eda-in-golang/internal/am"
)

type rawMessage struct {
	msg am.RawMessage
}

var _ am.RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string          { return m.msg.ID() }
func (m rawMessage) Subject() string     { return m.msg.Subject() }
func (m rawMessage) MessageName() string { return m.msg.MessageName() }
func (m rawMessage) Data() []byte        { return m.msg.Data() }
func (m *rawMessage) Ack() error         { return nil }
func (m *rawMessage) NAck() error        { return nil }
func (m rawMessage) Extend() error       { return nil }
func (m *rawMessage) Kill() error        { return nil }
