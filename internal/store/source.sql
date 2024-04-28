-- name: CreateGHSource :one
INSERT INTO gh_sources (
    workspace_id, owner, repo, path, installation_id
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListGHSources :many
SELECT * FROM gh_sources
WHERE workspace_id = $1;

-- name: GetGHSource :one
SELECT * FROM gh_sources
WHERE workspace_id = $1 AND id = $2;

-- name: DeleteGHSource :exec
DELETE FROM gh_sources
WHERE workspace_id = $1 AND id = $2;