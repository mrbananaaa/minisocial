package domain

import "github.com/google/uuid"

type PostCreated struct {
	PostID   uuid.UUID
	AuthorID uuid.UUID
	Title    string
}

func (PostCreated) IsEvent() {}
