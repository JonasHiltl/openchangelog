-- name: saveWorkspace :one
INSERT INTO workspaces (
    id, name
) VALUES (?1, ?2)
ON CONFLICT (id)
DO UPDATE SET name = ?2
RETURNING *;

-- name: getWorkspace :one
SELECT sqlc.embed(w), sqlc.embed(t)
FROM workspaces w
LEFT JOIN tokens t ON w.id = t.workspace_id
WHERE id = ?;

-- name: deleteWorkspace :exec
DELETE FROM workspaces
WHERE id = ?;

-- name: createToken :exec
INSERT INTO tokens (
    key, workspace_id
) VALUES (
    ?, ?
);

-- name: getToken :one
SELECT * FROM tokens
WHERE key = ?;

-- technically a workspace can have multiple tokens, but domain only create one token when workspace is created

-- name: getTokenByWorkspace :one
SELECT * FROM tokens
WHERE workspace_id = ?
LIMIT 1;

-- name: createChangelog :one
INSERT INTO changelogs (
    workspace_id, id, subdomain, title, subtitle, logo_src, logo_link, logo_alt, logo_height, logo_width
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: deleteChangelog :exec
DELETE FROM changelogs
WHERE workspace_id = ? AND id = ?;

-- name: getChangelog :one
SELECT sqlc.embed(c), sqlc.embed(cs)
FROM changelogs c
LEFT JOIN changelog_source cs ON c.workspace_id = cs.workspace_id AND c.source_id = cs.id
WHERE c.workspace_id = ? AND c.id = ?;

-- name: getChangelogBySubdomain :one
SELECT sqlc.embed(c), sqlc.embed(cs)
FROM changelogs c
LEFT JOIN changelog_source cs ON c.workspace_id = cs.workspace_id AND c.source_id = cs.id
WHERE c.subdomain = ?;

-- name: listChangelogs :many
SELECT sqlc.embed(c), sqlc.embed(cs)
FROM changelogs c
LEFT JOIN changelog_source cs ON c.workspace_id = cs.workspace_id AND c.source_id = cs.id
WHERE c.workspace_id = ?;

-- name: updateChangelog :one
UPDATE changelogs
SET
   title = coalesce(sqlc.narg(title), title),
   subtitle = coalesce(sqlc.narg(subtitle), subtitle),
   subdomain = coalesce(sqlc.narg(subdomain), subdomain),
   logo_src = coalesce(sqlc.narg(logo_src), logo_src),
   logo_link = coalesce(sqlc.narg(logo_link), logo_link),
   logo_alt = coalesce(sqlc.narg(logo_alt), logo_alt),
   logo_height = coalesce(sqlc.narg(logo_height), logo_height),
   logo_width = coalesce(sqlc.narg(logo_width), logo_width)
WHERE workspace_id = sqlc.arg(workspace_id) AND id = sqlc.arg(id)
RETURNING *;

-- name: setChangelogSource :exec
UPDATE changelogs
SET source_id = ?
WHERE workspace_id = ? AND id = ?;

-- name: deleteChangelogSource :exec
UPDATE changelogs
SET source_id = NULL
WHERE workspace_id = ? AND id = ?;

-- name: createGHSource :one
INSERT INTO gh_sources (
    id, workspace_id, owner, repo, path, installation_id
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: listGHSources :many
SELECT * FROM gh_sources
WHERE workspace_id = ?;

-- name: getGHSource :one
SELECT * FROM gh_sources
WHERE workspace_id = ? AND id = ?;

-- name: deleteGHSource :exec
DELETE FROM gh_sources
WHERE workspace_id = ? AND id = ?;