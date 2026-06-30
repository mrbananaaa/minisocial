package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres/sqlc"
)

type Repository struct {
	q *sqlc.Queries
}

func New(db sqlc.DBTX) *Repository {
	return &Repository{
		q: sqlc.New(db),
	}
}

func (r *Repository) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	row, err := r.q.CreatePost(ctx, fromDomain[sqlc.CreatePostParams](post))
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}

func (r *Repository) Update(ctx context.Context, post *domain.Post) error {
	_, err := r.q.UpdatePost(ctx, fromDomain[sqlc.UpdatePostParams](post))
	return mapError(err)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	row, err := r.q.FindPostByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}
