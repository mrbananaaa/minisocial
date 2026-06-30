package application

import (
	"context"

	"github.com/google/uuid"
)

func (a *Application) ArchivePost(ctx context.Context, id uuid.UUID) error {
	p, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := p.Archive(); err != nil {
		return err
	}

	err = a.repo.Update(ctx, p)
	if err != nil {
		return err
	}

	return nil
}
