package am

type (
	RawMessage interface {
		Message
		Data() []byte
	}

	rawMessage struct {
		id   string
		name string
		data []byte
	}
)

var _ RawMessage = (*rawMessage)(nil)

func (m rawMessage) ID() string          { return m.id }
func (m rawMessage) MessageName() string { return m.name }
func (m rawMessage) Data() []byte        { return m.data }
func (m rawMessage) Ack() error          { return nil }
func (m rawMessage) NAck() error         { return nil }
func (m rawMessage) Extend() error       { return nil }
func (m rawMessage) Kill() error         { return nil }
