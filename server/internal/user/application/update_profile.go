package application

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

type UpdateProfileInput struct {
	ID        uuid.UUID
	Name      string
	Bio       string
	AvatarURL string
}

func (a *Application) UpdateProfile(ctx context.Context, input UpdateProfileInput) (*domain.User, error) {
	user, err := a.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Name) != "" {
		if err := user.UpdateName(input.Name); err != nil {
			return nil, err
		}
	}

	if strings.TrimSpace(input.Bio) != "" {
		if err := user.UpdateBio(input.Bio); err != nil {
			return nil, err
		}
	}

	if strings.TrimSpace(input.AvatarURL) != "" {
		if err := user.UpdateAvatarURL(input.AvatarURL); err != nil {
			return nil, err
		}
	}

	if err := a.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
