package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventType string

const (
	EvPostCreated EventType = "post.created"
)

type PostCreated struct {
	PostID   uuid.UUID
	AuthorID uuid.UUID
	Title    string
}

func (PostCreated) EventType() string {
	return string(EvPostCreated)
}

func (PostCreated) AggregateType() string {
	return "post"
}

func (p PostCreated) AggregateID() uuid.UUID {
	return p.PostID
}

func (p PostCreated) Payload() []byte {
	payload, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}

	return payload
}
