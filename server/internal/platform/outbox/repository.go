package outbox

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/minisocial/internal/platform/db"
	"github.com/mrbananaaa/minisocial/internal/platform/outbox/sqlc"
)

type outboxRepo struct {
	q *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *outboxRepo {
	return &outboxRepo{
		q: sqlc.New(pool),
	}
}

// query returns transaction-bound query when the context carries
// a transaction; otherwise it returns the default query.
func (r *outboxRepo) query(ctx context.Context) *sqlc.Queries {
	if tx, ok := db.TxFromContext(ctx); ok {
		return r.q.WithTx(tx)
	}

	return r.q
}

func (r *outboxRepo) InsertEvent(ctx context.Context, event *OutboxEvents) error {
	q := r.query(ctx)

	return q.CreateOutboxEvents(ctx, sqlc.CreateOutboxEventsParams{
		ID:            event.ID,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID,
		EventType:     event.EventType,
		Payload:       event.Payload,
		CreatedAt:     event.CreatedAt,
	})
}
