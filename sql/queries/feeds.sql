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
SELECT
  feeds.name AS feed_name,
  feeds.url AS feed_url,
  users.name AS user_name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET
  last_fetched_at = NOW() AT TIME ZONE 'UTC',
  updated_at = NOW() AT TIME ZONE 'UTC'
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;