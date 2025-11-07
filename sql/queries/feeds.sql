-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :one
WITH new_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    nf.id,
    nf.created_at,
    nf.updated_at,
    nf.user_id,
    nf.feed_id,
    u.name AS user_name,
    f.name AS feed_name,
    f.url AS feed_url
FROM new_follow nf
JOIN users u ON nf.user_id = u.id
JOIN feeds f ON nf.feed_id = f.id;

-- name: GetFeedFromUrl :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: getFeedFollowsByUser :many
SELECT f.name,f.id,u.name FROM users u
    INNER JOIN feed_follows ff ON u.id = ff.user_id
    INNER JOIN feeds f ON ff.feed_id = f.id
    WHERE u.id = $1;
