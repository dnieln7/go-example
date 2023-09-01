-- name: CreatePost :one

INSERT INTO
    tb_posts(
        id,
        title,
        description,
        url,
        published_at,
        feed_id,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetPostForUser :many
SELECT tb_posts.* FROM tb_posts 
JOIN tb_feed_follows ON tb_posts.feed_id = tb_feed_follows.feed_id
WHERE tb_feed_follows.user_id = $1
ORDER BY tb_posts.published_at DESC
LIMIT $2;
