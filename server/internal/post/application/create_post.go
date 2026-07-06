package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

type CreatePostInput struct {
	AuthorID uuid.UUID
	Title    string
	Content  string
}

type CreatePostResult struct {
	Post   *domain.Post
	Events []events.Event
}

func (a *Application) CreatePost(
	ctx context.Context,
	input CreatePostInput,
) (*domain.Post, error) {
	p, err := domain.New(domain.NewPostInput(input))
	if err != nil {
		return nil, err
	}

	if err := a.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
