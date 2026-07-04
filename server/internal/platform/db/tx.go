package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/minisocial/internal/platform/events"
)

type txKey struct{}

type TxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) *TxManager {
	return &TxManager{
		db: db,
	}
}

func withTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// WithTx ...
// TODO: gave it options for BeginTx option
func (m *TxManager) WithTx(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	// Check if WithTx already inside a transaction
	if _, ok := TxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// INFO: use collector for testing or example, remove this latter
	collector := events.NewCollector()
	txCtx := withTx(ctx, tx)
	txCtx = events.ContextWithCollector(txCtx, collector)
	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
