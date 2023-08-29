-- +goose Up

CREATE TABLE tb_users (
    id UUID PRIMARY KEY,
    name TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down

DROP TABLE tb_users;
