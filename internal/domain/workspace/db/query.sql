-- name: SaveWorkspace :exec
INSERT INTO workspaces (
    id, name
) VALUES ($1, $2)
ON CONFLICT (id)
DO UPDATE SET name = $2;

-- name: GetWorkspace :one
SELECT sqlc.embed(w), sqlc.embed(t)
FROM workspaces w
LEFT JOIN tokens t ON w.id = t.workspace_id
WHERE id = $1;

-- name: DeleteWorkspace :exec
DELETE FROM workspaces
WHERE id = $1;

-- name: CreateToken :exec
INSERT INTO tokens (
    key, workspace_id
) VALUES (
    $1, $2
);

-- name: GetToken :one
SELECT * FROM tokens
WHERE key = $1;

-- technically a workspace can have multiple tokens, but domain only create one token when workspace is created

-- name: GetTokenByWorkspace :one
SELECT * FROM tokens
WHERE workspace_id = $1
LIMIT 1;

