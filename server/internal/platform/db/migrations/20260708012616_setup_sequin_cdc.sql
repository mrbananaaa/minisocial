-- +goose NO TRANSACTION
-- +goose Up
CREATE PUBLICATION sequin_pub FOR ALL TABLES WITH (publish_via_partition_root = true);

SELECT pg_create_logical_replication_slot('sequin_slot', 'pgoutput');

-- +goose NO TRANSACTION
-- +goose Down
SELECT pg_drop_replication_slot('sequin_slot');

DROP PUBLICATION sequin_pub;
