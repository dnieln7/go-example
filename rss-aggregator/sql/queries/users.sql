-- name: CreateUser :one
INSERT INTO tb_users 
(id, created_at, updated_at, name, api_key) 
VALUES 
($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex')) 
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM tb_users WHERE api_key = $1;