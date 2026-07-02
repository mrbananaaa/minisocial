package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/db"
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

// query returns transaction-bound query when the context carries
// a transaction; otherwise it returns the default query.
func (r *Repository) query(ctx context.Context) *sqlc.Queries {
	if tx, ok := db.TxFromContext(ctx); ok {
		return r.q.WithTx(tx)
	}

	return r.q
}

func (r *Repository) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	q := r.query(ctx)

	row, err := q.CreatePost(ctx, fromDomain[sqlc.CreatePostParams](post))
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}

func (r *Repository) Update(ctx context.Context, post *domain.Post) error {
	q := r.query(ctx)

	_, err := q.UpdatePost(ctx, fromDomain[sqlc.UpdatePostParams](post))
	return mapError(err)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	q := r.query(ctx)

	row, err := q.FindPostByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	return toDomain(row), nil
}
