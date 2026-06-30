package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/minisocial/internal/platform/config"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DatabaseConfig) (*Database, error) {
	pool, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open pgx pool: %v", err)
	}

	db := &Database{
		pool: pool,
	}

	if err := db.Health(ctx); err != nil {
		return nil, fmt.Errorf("Couldn't reach database: %v", err)
	}

	return db, nil
}

func (db *Database) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *Database) Close() {
	db.pool.Close()
}
