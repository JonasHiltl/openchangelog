-- name: CreateChangelog :one
INSERT INTO changelogs (
    workspace_id, title, subtitle, logo_src, logo_link, logo_alt, logo_height, logo_width
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetChangelog :one
SELECT * FROM changelogs
WHERE workspace_id = $1 AND id =$2;

-- name: DeleteChangelog :exec
DELETE FROM changelogs
WHERE workspace_id = $1 AND id =$2;

-- name: ListChangelogs :many
SELECT * FROM changelogs
WHERE workspace_id = $1;

-- name: UpdateChangelog :one
UPDATE changelogs
SET
   title = coalesce(sqlc.narg(title), title),
   subtitle = coalesce(sqlc.narg(subtitle), subtitle),
   logo_src = coalesce(sqlc.narg(logo_src), logo_src),
   logo_link = coalesce(sqlc.narg(logo_link), logo_link),
   logo_alt = coalesce(sqlc.narg(logo_alt), logo_alt),
   logo_height = coalesce(sqlc.narg(logo_height), logo_height),
   logo_width = coalesce(sqlc.narg(logo_width), logo_width)
WHERE workspace_id = sqlc.arg(workspace_id) AND id = sqlc.arg(id)
RETURNING *;

-- name: UpdateChangelogSource :exec
UPDATE changelogs 
SET source_id = $1, source_type = $2
WHERE workspace_id = $3 AND id = $4;

CREATE VIEW changelog_sources AS (
  SELECT gh.* FROM changelogs
  LEFT JOIN gh_sources gh ON cl.source_type = 'GitHub' AND gh.id = cl.source_id
);

-- name: GetChangelogSource :one
SELECT sqlc.embed(cl), sqlc.embed(s)
FROM changelogs cl
LEFT JOIN changelog_sources s on cl.source_id = s.id
WHERE cl.workspace_id = $1 AND cl.id = $2;