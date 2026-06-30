package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*Post, error)
	// GetBySlug(ctx context.Context, slug string) (*Post, error)
}
