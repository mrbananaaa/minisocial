package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationsFS embed.FS

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	dbconn := stdlib.OpenDBFromPool(pool)
	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		return fmt.Errorf("Couldn't reach database: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("Failed to select dialect: %v", err)
	}
	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(goose.NopLogger())

	if err := goose.UpContext(ctx, dbconn, "."); err != nil {
		return fmt.Errorf("Failed to migrate database: %v", err)
	}

	return nil
}
