package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OutboxEvents struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   uuid.UUID
	EventType     string
	Payload       []byte
	CreatedAt     time.Time
}

func NewEvent(
	aggregateType string,
	aggregateID uuid.UUID,
	eventType string,
	payload []byte,
) *OutboxEvents {
	return &OutboxEvents{
		ID:            uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payload,
		CreatedAt:     time.Now(),
	}
}

type Repository interface {
	InsertEvent(context.Context, *OutboxEvents) error
}
