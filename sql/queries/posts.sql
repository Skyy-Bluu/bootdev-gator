-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUserByUserID :many
SELECT * FROM posts 
WHERE feed_id in (
    SELECT feeds.id FROM feed_follows 
    JOIN users ON users.id = feed_follows.user_id AND feed_follows.user_id = $1
    JOIN feeds ON feeds.id = feed_follows.feed_id
)
ORDER BY published_at
LIMIT $2;

