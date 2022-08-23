package am

type (
	RawMessageStream     = MessageStream[RawMessage, IncomingRawMessage]
	RawMessageHandler    = MessageHandler[IncomingRawMessage]
	RawMessagePublisher  = MessagePublisher[RawMessage]
	RawMessageSubscriber = MessageSubscriber[IncomingRawMessage]

	RawMessage interface {
		Message
		Data() []byte
	}

	IncomingRawMessage interface {
		IncomingMessage
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
