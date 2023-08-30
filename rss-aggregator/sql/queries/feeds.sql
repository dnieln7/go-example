-- name: CreateFeed :one
INSERT INTO tb_feeds
(id, created_at, updated_at, name, url, user_id) 
VALUES 
($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM tb_feeds;