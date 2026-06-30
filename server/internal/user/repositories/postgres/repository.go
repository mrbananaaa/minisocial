package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
	"github.com/mrbananaaa/minisocial/internal/user/repositories/postgres/sqlc"
)

type Repository struct {
	q *sqlc.Queries
}

func New(db sqlc.DBTX) *Repository {
	return &Repository{
		q: sqlc.New(db),
	}
}

func (r *Repository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	row, err := r.q.CreateUser(ctx, fromDomain(user))
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.q.FindUserByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := r.q.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}

func (r *Repository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	row, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        user.ID,
		Name:      user.Name,
		Bio:       &user.Bio,
		AvatarUrl: &user.AvatarURL,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteUser(ctx, id); err != nil {
		return mapError(err)
	}

	return nil
}
