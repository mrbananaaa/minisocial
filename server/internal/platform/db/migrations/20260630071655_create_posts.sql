-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY NOT NULL,
  author_id UUID NOT NULL,
  title TEXT NOT NULL,
  slug TEXT NOT NULL,
  content TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  archived_at TIMESTAMPTZ NULL,

  CONSTRAINT fk_user_post
    FOREIGN KEY (author_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS posts_slug_idx
  ON posts(slug);
CREATE INDEX IF NOT EXISTS posts_author_id_idx
  ON posts(author_id);

-- +goose Down
DROP INDEX IF EXISTS posts_slug_idx;
DROP INDEX IF EXISTS posts_author_id_idx;

DROP TABLE IF EXISTS posts;
