-- +goose Up
CREATE TABLE IF NOT EXISTS outbox_events (
  id UUID PRIMARY KEY NOT NULL,
  aggregate_type TEXT NOT NULL,
  aggregate_id UUID NOT NULL,
  event_type TEXT NOT NULL,
  payload BYTEA NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_outbox_events_created_at
  ON outbox_events(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_outbox_events_created_at;

DROP TABLE IF EXISTS outbox_events;
