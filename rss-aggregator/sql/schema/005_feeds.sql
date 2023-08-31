-- +goose Up

ALTER TABLE tb_feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down

ALTER TABLE tb_feeds DROP COLUMN last_fetched_at;
