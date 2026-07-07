package messaging

type Envelope struct {
	ID      string
	Topic   string
	Payload []byte

	ack func() error
}

func NewEnvelope(
	id string,
	topic string,
	payload []byte,
	ack func() error,
) Envelope {

	return Envelope{
		ID:      id,
		Topic:   topic,
		Payload: payload,
		ack:     ack,
	}
}

func (e Envelope) Ack() error {
	if e.ack == nil {
		return nil
	}

	return e.ack()
}

type Subscription struct {
	messages <-chan Envelope
	errors   <-chan error
	close    func()
}

func NewSubscription(
	messages <-chan Envelope,
	errors <-chan error,
	closeFn func(),
) *Subscription {

	return &Subscription{
		messages: messages,
		errors:   errors,
		close:    closeFn,
	}
}

func (s *Subscription) Messages() <-chan Envelope {
	return s.messages
}

func (s *Subscription) Errors() <-chan error {
	return s.errors
}

func (s *Subscription) Close() {
	if s.close != nil {
		s.close()
	}
}
