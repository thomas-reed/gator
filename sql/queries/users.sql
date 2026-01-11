-- name: CreateUser :one
INSERT INTO users (id, name, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  $1,
  NOW() AT TIME ZONE 'UTC',
  NOW() AT TIME ZONE 'UTC'
)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name FROM users;