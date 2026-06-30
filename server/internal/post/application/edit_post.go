package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

type EditPostInput struct {
	ID      uuid.UUID
	Title   *string
	Content *string
}

func (a *Application) EditPost(ctx context.Context, input EditPostInput) (*domain.Post, error) {
	p, err := a.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	var ut, uc = p.Title, p.Content

	if input.Title != nil {
		ut = *input.Title
	}

	if input.Content != nil {
		uc = *input.Content
	}

	if err := p.Edit(ut, uc); err != nil {
		return nil, err
	}

	err = a.repo.Update(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
