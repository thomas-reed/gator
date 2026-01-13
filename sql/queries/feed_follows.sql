-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
  VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC'
  )
  RETURNING *
)
SELECT
  inserted_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsByUser :many
SELECT
  users.name AS user_name,
  feeds.name AS feed_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE users.id = $1;

-- name: DeleteFeedFollowByUserAndName :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;