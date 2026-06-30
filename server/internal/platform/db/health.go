package db

import (
	"context"
)

func (db *Database) Health(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
