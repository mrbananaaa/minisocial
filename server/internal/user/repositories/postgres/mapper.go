package postgres

import (
	"github.com/mrbananaaa/minisocial/internal/user/domain"
	"github.com/mrbananaaa/minisocial/internal/user/repositories/postgres/sqlc"
)

func toDomain(u sqlc.User) *domain.User {
	return &domain.User{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		Bio:          *u.Bio,
		AvatarURL:    *u.AvatarUrl,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func fromDomain(u *domain.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		Bio:          &u.Bio,
		AvatarUrl:    &u.AvatarURL,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
