package messaging

type Message struct {
	ID      string
	Topic   string
	Payload []byte
	AckFunc func() error
}

func (m *Message) Ack() error {
	if m.AckFunc != nil {
		m.AckFunc()
	}
	return nil
}

type SubscribtionPayload struct {
	Event   <-chan Message
	ErrChan <-chan error
}
