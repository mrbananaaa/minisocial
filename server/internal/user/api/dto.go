package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/user/application"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name,omitempty" validate:"min=5,max=50"`
	Bio       string `json:"bio,omitempty" validate:"min=5,max=255"`
	AvatarURL string `json:"avatar_url,omitempty" validate:"min=3,max=255"`
}

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
