-- name: CreateOutboxEvents :exec
INSERT INTO outbox_events (
  id,
  aggregate_type,
  aggregate_id,
  event_type,
  payload,
  created_at
) VALUES (
  $1, $2, $3, $4, $5, $6
);
