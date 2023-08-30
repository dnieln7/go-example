-- +goose Up

ALTER TABLE tb_users ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
    encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down

ALTER TABLE tb_users DROP COLUMN api_key;
