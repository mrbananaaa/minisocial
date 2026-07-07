package domain

import "github.com/google/uuid"

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
