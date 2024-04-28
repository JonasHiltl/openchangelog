-- name: CreateToken :one
INSERT INTO tokens (
    key, workspace_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetToken :one
SELECT * FROM tokens
WHERE key = $1;

-- name: ListToken :many
SELECT * FROM tokens
WHERE workspace_id = $1;