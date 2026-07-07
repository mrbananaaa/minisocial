package events

import (
	"context"

	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
)

type Event interface {
	EventType() string
}

type EventMessage struct {
	Topic   string
	Payload []byte
}

type EventSource interface {
	PullEvents() []Event
}

type EventHandler interface {
	EventType() string
	Handle(context.Context, EventMessage) error
}

type Consumer interface {
	Subscribe(ctx context.Context, topic string) (*messaging.Subscription, error)
}
