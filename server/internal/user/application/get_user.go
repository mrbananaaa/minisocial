package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

func (a *Application) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *Application) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return a.repo.GetByEmail(ctx, email)
}

func (a *Application) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return a.repo.GetByUsername(ctx, username)
}

type GetUserInput struct {
	ID       *uuid.UUID
	Email    *string
	Username *string
}

func (a *Application) GetUser(ctx context.Context, input GetUserInput) (*domain.User, error) {
	if input.ID != nil {
		u, err := a.repo.GetByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	if input.Email != nil {
		u, err := a.repo.GetByEmail(ctx, *input.Email)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	if input.Username != nil {
		u, err := a.repo.GetByUsername(ctx, *input.Username)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	return nil, nil
}
