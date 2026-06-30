-- name: CreateUser :one
INSERT INTO users (
  id,
  email,
  username,
  name,
  password_hash,
  bio,
  avatar_url,
  created_at,
  updated_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9
) RETURNING *;

-- name: FindUserByID :one
SELECT
  id,
  email,
  username,
  name,
  password_hash,
  bio,
  avatar_url,
  created_at,
  updated_at
FROM users
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT
  id,
  email,
  username,
  name,
  password_hash,
  bio,
  avatar_url,
  created_at,
  updated_at
FROM users
WHERE email = $1;

-- name: FindUserByUsername :one
SELECT
  id,
  email,
  username,
  name,
  password_hash,
  bio,
  avatar_url,
  created_at,
  updated_at
FROM users
WHERE username = $1;

-- name: UpdateUser :one
UPDATE users
SET
  name = $2,
  bio = $3,
  avatar_url = $4,
  updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;