package api

import (
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/user/application"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

func toCreateUserInput(req CreateUserRequest) application.CreateUserInput {
	return application.CreateUserInput{
		Email:    req.Email,
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	}
}

func toUpdateUserInput(id uuid.UUID, req UpdateProfileRequest) application.UpdateProfileInput {
	return application.UpdateProfileInput{
		ID:        id,
		Name:      req.Name,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
	}
}

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
