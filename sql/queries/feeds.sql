-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  $1,
  $2, 
  $3,
  NOW() AT TIME ZONE 'UTC',
  NOW() AT TIME ZONE 'UTC'
)
RETURNING *;

-- name: ListFeeds :many
SELECT feeds.name, feeds.url, users.name FROM feeds
INNER JOIN users ON feeds.user_id = users.id;