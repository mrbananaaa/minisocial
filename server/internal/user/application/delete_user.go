package application

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, user.ID)
}
