package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

type CreatePostInput struct {
	AuthorID uuid.UUID
	Title    string
	Content  string
}

func (a *Application) CreatePost(ctx context.Context, input CreatePostInput) (*domain.Post, error) {
	// TODO: check if author id is valid user

	p, err := domain.New(domain.NewPostInput(input))
	if err != nil {
		return nil, err
	}

	post, err := a.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return post, nil
}
