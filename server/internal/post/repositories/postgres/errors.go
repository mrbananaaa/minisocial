package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

const (
	slugUniqueConstraint = "posts_slug_idx"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrPostNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case slugUniqueConstraint:
			return domain.ErrPostSlugAlreadyExists

		default:
			return fmt.Errorf("unknown user repository error: %w", err)
		}
	}

	return err
}
