package api

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
