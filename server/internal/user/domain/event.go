package domain

import "github.com/google/uuid"

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
