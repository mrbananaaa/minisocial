package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

type GetUserInput struct {
	ID       *uuid.UUID
	Email    *string
	Username *string
}

func (s *Service) GetUser(ctx context.Context, input GetUserInput) (*domain.User, error) {
	if input.ID != nil {
		u, err := s.repo.GetByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	if input.Email != nil {
		u, err := s.repo.GetByEmail(ctx, *input.Email)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	if input.Username != nil {
		u, err := s.repo.GetByUsername(ctx, *input.Username)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	return nil, nil
}
