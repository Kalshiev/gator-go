-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListFeeds :many
SELECT *
FROM feeds;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: CreateFeedFollow :one
WITH inserted_new_follow AS (
    INSERT INTO feeds_follow (id, created_at, updated_at, feed_id, user_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
inserted_new_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM inserted_new_follow
INNER JOIN feeds ON inserted_new_follow.feed_id = feeds.id
INNER JOIN users ON inserted_new_follow.user_id = users.id;

-- name: GetFeedFollowForUser :many
SELECT
feeds_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM feeds_follow
INNER JOIN feeds ON feeds_follow.feed_id = feeds.id
INNER JOIN users ON feeds_follow.user_id = users.id
WHERE feeds_follow.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds_follow
WHERE feed_id = $1 AND user_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1,
    updated_at = $1
WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;