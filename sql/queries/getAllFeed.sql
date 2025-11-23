-- name: GetAllFeed :many
SELECT * FROM feeds;

-- name: GetFeedFromURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedFromID :one
SELECT * FROM feeds WHERE id = $1;
