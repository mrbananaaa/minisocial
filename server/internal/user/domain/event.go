package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventType string

const (
	EvUserCreated EventType = "user.created"
)

type UserCreated struct {
	UserID uuid.UUID
	Email  string
}

func (UserCreated) EventType() string {
	return string(EvUserCreated)
}

func (UserCreated) AggregateType() string {
	return "post"
}

func (p UserCreated) AggregateID() uuid.UUID {
	return p.UserID
}

func (p UserCreated) Payload() []byte {
	payload, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}

	return payload
}
