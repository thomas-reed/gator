-- name: CreatePost :one
INSERT INTO posts (id, title, url, published_at, description, feed_id, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  $4,
  $5,
  NOW() AT TIME ZONE 'UTC',
  NOW() AT TIME ZONE 'UTC'
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name  FROM posts
INNER JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
INNER JOIN feeds ON posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY published_at DESC NULLS LAST
LIMIT $2;

-- name: GetPostByURL :one
SELECT * FROM posts
WHERE url = $1;
