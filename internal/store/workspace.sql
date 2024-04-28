-- name: CreateWorkspace :one
INSERT INTO workspaces (
    id, name
) VALUES ($1, $2)
RETURNING *;

-- name: GetWorkspace :one
SELECT sqlc.embed(w), sqlc.embed(t) 
FROM workspaces as w
JOIN tokens as t ON t.workspace_id = w.id
WHERE id = $1;

-- name: UpdateWorkspace :one
UPDATE workspaces
SET name = coalesce(sqlc.narg(name), name)
WHERE id = $1
RETURNING *;

-- name: DeleteWorkspace :exec
DELETE FROM workspaces
WHERE id =$1;