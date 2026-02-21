-- name: CreateFeedFollows :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *, (SELECT name AS user_name FROM users WHERE id = $4) , (SELECT name AS feed_name FROM feeds WHERE id = $5);

-- name: GetFeedFollowsByUser :many
SELECT users.name AS user_name, feeds.name AS feed_name FROM feed_follows 
JOIN users ON users.id = feed_follows.user_id AND feed_follows.user_id = $1
JOIN feeds ON feeds.id = feed_follows.feed_id;

-- name: DeleteFeedFollowEntryByUserIDAndFeedID :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;
