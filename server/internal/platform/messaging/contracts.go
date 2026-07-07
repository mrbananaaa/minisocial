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
	Message <-chan Message
	Errors  <-chan error
	Close   func()
}
