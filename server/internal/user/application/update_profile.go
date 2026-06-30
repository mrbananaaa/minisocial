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
		user.Name = input.Name
	}

	if strings.TrimSpace(input.Bio) != "" {
		user.Bio = input.Bio
	}

	if strings.TrimSpace(input.AvatarURL) != "" {
		user.AvatarURL = input.AvatarURL
	}

	updatedUser, err := a.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
