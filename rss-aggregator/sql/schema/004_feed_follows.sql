-- +goose Up

CREATE TABLE tb_feed_follows (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES tb_users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES tb_feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(user_id, feed_id)
);

-- +goose Down

DROP TABLE tb_feed_follows;
