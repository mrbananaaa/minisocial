package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

const (
	emailUniqueConstraint    = "users_email_idx"
	usernameUniqueConstraint = "users_username_idx"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case emailUniqueConstraint:
			return domain.ErrEmailAlreadyExists

		case usernameUniqueConstraint:
			return domain.ErrUsernameAlreadyExists

		default:
			return fmt.Errorf("unknown user repository error: %w", err)
		}
	}

	return err
}
