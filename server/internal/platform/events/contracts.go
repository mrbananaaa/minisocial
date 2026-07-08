package events

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
)

type Event interface {
	AggregateType() string
	AggregateID() uuid.UUID
	EventType() string
	Payload() []byte
}

// TODO: unused for now

type EventMessage struct {
	Topic   string
	Payload []byte
}

// INFO: it's just a wrapper in case we need more fields

type EventContext struct {
	Envelope messaging.Envelope
	EventID  string
}

type EventSource interface {
	PullEvents() []Event
}

type EventHandler interface {
	EventType() string
	Handle(context.Context, EventContext) error
}

type Consumer interface {
	Subscribe(ctx context.Context, topic string) (*messaging.Subscription, error)
}
