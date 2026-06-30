package application

import (
	"context"

	"github.com/google/uuid"
)

func (a *Application) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return a.repo.Delete(ctx, user.ID)
}
