-- +goose Up

CREATE TABLE tb_feeds (
    id UUID PRIMARY KEY,
    name TEXT,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES tb_users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down

DROP TABLE tb_feeds;
