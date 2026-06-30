-- name: CreatePost :one
INSERT INTO posts (
  id,
  author_id,
  title,
  slug,
  content,
  status,
  created_at,
  updated_at,
  archived_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9
) RETURNING *;

-- name: UpdatePost :one
UPDATE posts
SET
  title = $2,
  slug = $3,
  content = $4,
  status = $5,
  updated_at = $6,
  archived_at = $7
WHERE id = $1
RETURNING *;

-- name: FindPostByID :one
SELECT
  id,
  author_id,
  title,
  slug,
  content,
  status,
  created_at,
  updated_at,
  archived_at
FROM posts
WHERE id = $1;