// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package store

import (
	"context"
	"database/sql"
)

const createChangelog = `-- name: createChangelog :one
INSERT INTO changelogs (
    workspace_id, id, subdomain, domain, title, subtitle, logo_src, logo_link, logo_alt, logo_height, logo_width
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, workspace_id, subdomain, title, subtitle, source_id, logo_src, logo_link, logo_alt, logo_height, logo_width, created_at, domain
`

type createChangelogParams struct {
	WorkspaceID string
	ID          string
	Subdomain   string
	Domain      sql.NullString
	Title       sql.NullString
	Subtitle    sql.NullString
	LogoSrc     sql.NullString
	LogoLink    sql.NullString
	LogoAlt     sql.NullString
	LogoHeight  sql.NullString
	LogoWidth   sql.NullString
}

func (q *Queries) createChangelog(ctx context.Context, arg createChangelogParams) (changelog, error) {
	row := q.db.QueryRowContext(ctx, createChangelog,
		arg.WorkspaceID,
		arg.ID,
		arg.Subdomain,
		arg.Domain,
		arg.Title,
		arg.Subtitle,
		arg.LogoSrc,
		arg.LogoLink,
		arg.LogoAlt,
		arg.LogoHeight,
		arg.LogoWidth,
	)
	var i changelog
	err := row.Scan(
		&i.ID,
		&i.WorkspaceID,
		&i.Subdomain,
		&i.Title,
		&i.Subtitle,
		&i.SourceID,
		&i.LogoSrc,
		&i.LogoLink,
		&i.LogoAlt,
		&i.LogoHeight,
		&i.LogoWidth,
		&i.CreatedAt,
		&i.Domain,
	)
	return i, err
}

const createGHSource = `-- name: createGHSource :one
INSERT INTO gh_sources (
    id, workspace_id, owner, repo, path, installation_id
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, workspace_id, owner, repo, path, installation_id
`

type createGHSourceParams struct {
	ID             string
	WorkspaceID    string
	Owner          string
	Repo           string
	Path           string
	InstallationID int64
}

func (q *Queries) createGHSource(ctx context.Context, arg createGHSourceParams) (ghSource, error) {
	row := q.db.QueryRowContext(ctx, createGHSource,
		arg.ID,
		arg.WorkspaceID,
		arg.Owner,
		arg.Repo,
		arg.Path,
		arg.InstallationID,
	)
	var i ghSource
	err := row.Scan(
		&i.ID,
		&i.WorkspaceID,
		&i.Owner,
		&i.Repo,
		&i.Path,
		&i.InstallationID,
	)
	return i, err
}

const createToken = `-- name: createToken :exec
INSERT INTO tokens (
    key, workspace_id
) VALUES (
    ?, ?
)
`

type createTokenParams struct {
	Key         string
	WorkspaceID string
}

func (q *Queries) createToken(ctx context.Context, arg createTokenParams) error {
	_, err := q.db.ExecContext(ctx, createToken, arg.Key, arg.WorkspaceID)
	return err
}

const deleteChangelog = `-- name: deleteChangelog :exec
DELETE FROM changelogs
WHERE workspace_id = ? AND id = ?
`

type deleteChangelogParams struct {
	WorkspaceID string
	ID          string
}

func (q *Queries) deleteChangelog(ctx context.Context, arg deleteChangelogParams) error {
	_, err := q.db.ExecContext(ctx, deleteChangelog, arg.WorkspaceID, arg.ID)
	return err
}

const deleteChangelogSource = `-- name: deleteChangelogSource :exec
UPDATE changelogs
SET source_id = NULL
WHERE workspace_id = ? AND id = ?
`

type deleteChangelogSourceParams struct {
	WorkspaceID string
	ID          string
}

func (q *Queries) deleteChangelogSource(ctx context.Context, arg deleteChangelogSourceParams) error {
	_, err := q.db.ExecContext(ctx, deleteChangelogSource, arg.WorkspaceID, arg.ID)
	return err
}

const deleteGHSource = `-- name: deleteGHSource :exec
DELETE FROM gh_sources
WHERE workspace_id = ? AND id = ?
`

type deleteGHSourceParams struct {
	WorkspaceID string
	ID          string
}

func (q *Queries) deleteGHSource(ctx context.Context, arg deleteGHSourceParams) error {
	_, err := q.db.ExecContext(ctx, deleteGHSource, arg.WorkspaceID, arg.ID)
	return err
}

const deleteWorkspace = `-- name: deleteWorkspace :exec
DELETE FROM workspaces
WHERE id = ?
`

func (q *Queries) deleteWorkspace(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteWorkspace, id)
	return err
}

const getChangelog = `-- name: getChangelog :one
SELECT c.id, c.workspace_id, c.subdomain, c.title, c.subtitle, c.source_id, c.logo_src, c.logo_link, c.logo_alt, c.logo_height, c.logo_width, c.created_at, c.domain, cs.id, cs.workspace_id, cs.owner, cs.repo, cs.path, cs.installation_id
FROM changelogs c
LEFT JOIN changelog_source cs ON c.workspace_id = cs.workspace_id AND c.source_id = cs.id
WHERE c.workspace_id = ? AND c.id = ?
`

type getChangelogParams struct {
	WorkspaceID string
	ID          string
}

type getChangelogRow struct {
	changelog       changelog
	ChangelogSource changelogSource
}

func (q *Queries) getChangelog(ctx context.Context, arg getChangelogParams) (getChangelogRow, error) {
	row := q.db.QueryRowContext(ctx, getChangelog, arg.WorkspaceID, arg.ID)
	var i getChangelogRow
	err := row.Scan(
		&i.changelog.ID,
		&i.changelog.WorkspaceID,
		&i.changelog.Subdomain,
		&i.changelog.Title,
		&i.changelog.Subtitle,
		&i.changelog.SourceID,
		&i.changelog.LogoSrc,
		&i.changelog.LogoLink,
		&i.changelog.LogoAlt,
		&i.changelog.LogoHeight,
		&i.changelog.LogoWidth,
		&i.changelog.CreatedAt,
		&i.changelog.Domain,
		&i.ChangelogSource.ID,
		&i.ChangelogSource.WorkspaceID,
		&i.ChangelogSource.Owner,
		&i.ChangelogSource.Repo,
		&i.ChangelogSource.Path,
		&i.ChangelogSource.InstallationID,
	)
	return i, err
}

const getChangelogByDomainOrSubdomain = `-- name: getChangelogByDomainOrSubdomain :one
SELECT c.id, c.workspace_id, c.subdomain, c.title, c.subtitle, c.source_id, c.logo_src, c.logo_link, c.logo_alt, c.logo_height, c.logo_width, c.created_at, c.domain, cs.id, cs.workspace_id, cs.owner, cs.repo, cs.path, cs.installation_id
FROM changelogs c
LEFT JOIN changelog_source cs ON c.workspace_id = cs.workspace_id AND c.source_id = cs.id
WHERE c.domain = ? OR c.subdomain = ?
LIMIT 1
`

type getChangelogByDomainOrSubdomainParams struct {
	Domain    sql.NullString
	Subdomain string
}

type getChangelogByDomainOrSubdomainRow struct {
	changelog       changelog
	ChangelogSource changelogSource
}

func (q *Queries) getChangelogByDomainOrSubdomain(ctx context.Context, arg getChangelogByDomainOrSubdomainParams) (getChangelogByDomainOrSubdomainRow, error) {
	row := q.db.QueryRowContext(ctx, getChangelogByDomainOrSubdomain, arg.Domain, arg.Subdomain)
	var i getChangelogByDomainOrSubdomainRow
	err := row.Scan(
		&i.changelog.ID,
		&i.changelog.WorkspaceID,
		&i.changelog.Subdomain,
		&i.changelog.Title,
		&i.changelog.Subtitle,
		&i.changelog.SourceID,
		&i.changelog.LogoSrc,
		&i.changelog.LogoLink,
		&i.changelog.LogoAlt,
		&i.changelog.LogoHeight,
		&i.changelog.LogoWidth,
		&i.changelog.CreatedAt,
		&i.changelog.Domain,
		&i.ChangelogSource.ID,
		&i.ChangelogSource.WorkspaceID,
		&i.ChangelogSource.Owner,
		&i.ChangelogSource.Repo,
		&i.ChangelogSource.Path,
		&i.ChangelogSource.InstallationID,
	)
	return i, err
}

const getGHSource = `-- name: getGHSource :one
SELECT id, workspace_id, owner, repo, path, installation_id FROM gh_sources
WHERE workspace_id = ? AND id = ?
`

type getGHSourceParams struct {
	WorkspaceID string
	ID          string
}

func (q *Queries) getGHSource(ctx context.Context, arg getGHSourceParams) (ghSource, error) {
	row := q.db.QueryRowContext(ctx, getGHSource, arg.WorkspaceID, arg.ID)
	var i ghSource
	err := row.Scan(
		&i.ID,
		&i.WorkspaceID,
		&i.Owner,
		&i.Repo,
		&i.Path,
		&i.InstallationID,
	)
	return i, err
}

const getToken = `-- name: getToken :one
SELECT "key", workspace_id FROM tokens
WHERE key = ?
`

func (q *Queries) getToken(ctx context.Context, key string) (token, error) {
	row := q.db.QueryRowContext(ctx, getToken, key)
	var i token
	err := row.Scan(&i.Key, &i.WorkspaceID)
	return i, err
}

const getTokenByWorkspace = `-- name: getTokenByWorkspace :one

SELECT "key", workspace_id FROM tokens
WHERE workspace_id = ?
LIMIT 1
`

// technically a workspace can have multiple tokens, but domain only create one token when workspace is created
func (q *Queries) getTokenByWorkspace(ctx context.Context, workspaceID string) (token, error) {
	row := q.db.QueryRowContext(ctx, getTokenByWorkspace, workspaceID)
	var i token
	err := row.Scan(&i.Key, &i.WorkspaceID)
	return i, err
}

const getWorkspace = `-- name: getWorkspace :one
SELECT w.id, w.name, t."key", t.workspace_id
FROM workspaces w
LEFT JOIN tokens t ON w.id = t.workspace_id
WHERE id = ?
`

type getWorkspaceRow struct {
	workspace workspace
	token     token
}

func (q *Queries) getWorkspace(ctx context.Context, id string) (getWorkspaceRow, error) {
	row := q.db.QueryRowContext(ctx, getWorkspace, id)
	var i getWorkspaceRow
	err := row.Scan(
		&i.workspace.ID,
		&i.workspace.Name,
		&i.token.Key,
		&i.token.WorkspaceID,
	)
	return i, err
}

const listChangelogs = `-- name: listChangelogs :many
SELECT c.id, c.workspace_id, c.subdomain, c.title, c.subtitle, c.source_id, c.logo_src, c.logo_link, c.logo_alt, c.logo_height, c.logo_width, c.created_at, c.domain, cs.id, cs.workspace_id, cs.owner, cs.repo, cs.path, cs.installation_id
FROM changelogs c
LEFT JOIN changelog_source cs ON c.workspace_id = cs.workspace_id AND c.source_id = cs.id
WHERE c.workspace_id = ?
`

type listChangelogsRow struct {
	changelog       changelog
	ChangelogSource changelogSource
}

func (q *Queries) listChangelogs(ctx context.Context, workspaceID string) ([]listChangelogsRow, error) {
	rows, err := q.db.QueryContext(ctx, listChangelogs, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []listChangelogsRow
	for rows.Next() {
		var i listChangelogsRow
		if err := rows.Scan(
			&i.changelog.ID,
			&i.changelog.WorkspaceID,
			&i.changelog.Subdomain,
			&i.changelog.Title,
			&i.changelog.Subtitle,
			&i.changelog.SourceID,
			&i.changelog.LogoSrc,
			&i.changelog.LogoLink,
			&i.changelog.LogoAlt,
			&i.changelog.LogoHeight,
			&i.changelog.LogoWidth,
			&i.changelog.CreatedAt,
			&i.changelog.Domain,
			&i.ChangelogSource.ID,
			&i.ChangelogSource.WorkspaceID,
			&i.ChangelogSource.Owner,
			&i.ChangelogSource.Repo,
			&i.ChangelogSource.Path,
			&i.ChangelogSource.InstallationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listGHSources = `-- name: listGHSources :many
SELECT id, workspace_id, owner, repo, path, installation_id FROM gh_sources
WHERE workspace_id = ?
`

func (q *Queries) listGHSources(ctx context.Context, workspaceID string) ([]ghSource, error) {
	rows, err := q.db.QueryContext(ctx, listGHSources, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ghSource
	for rows.Next() {
		var i ghSource
		if err := rows.Scan(
			&i.ID,
			&i.WorkspaceID,
			&i.Owner,
			&i.Repo,
			&i.Path,
			&i.InstallationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveWorkspace = `-- name: saveWorkspace :one
INSERT INTO workspaces (
    id, name
) VALUES (?1, ?2)
ON CONFLICT (id)
DO UPDATE SET name = ?2
RETURNING id, name
`

type saveWorkspaceParams struct {
	ID   string
	Name string
}

func (q *Queries) saveWorkspace(ctx context.Context, arg saveWorkspaceParams) (workspace, error) {
	row := q.db.QueryRowContext(ctx, saveWorkspace, arg.ID, arg.Name)
	var i workspace
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const setChangelogSource = `-- name: setChangelogSource :exec
UPDATE changelogs
SET source_id = ?
WHERE workspace_id = ? AND id = ?
`

type setChangelogSourceParams struct {
	SourceID    sql.NullString
	WorkspaceID string
	ID          string
}

func (q *Queries) setChangelogSource(ctx context.Context, arg setChangelogSourceParams) error {
	_, err := q.db.ExecContext(ctx, setChangelogSource, arg.SourceID, arg.WorkspaceID, arg.ID)
	return err
}

const updateChangelog = `-- name: updateChangelog :one
UPDATE changelogs
SET
   title = coalesce(?1, title),
   subtitle = coalesce(?2, subtitle),
   subdomain = coalesce(?3, subdomain),
   domain = coalesce(?4, domain),
   logo_src = coalesce(?5, logo_src),
   logo_link = coalesce(?6, logo_link),
   logo_alt = coalesce(?7, logo_alt),
   logo_height = coalesce(?8, logo_height),
   logo_width = coalesce(?9, logo_width)
WHERE workspace_id = ?10 AND id = ?11
RETURNING id, workspace_id, subdomain, title, subtitle, source_id, logo_src, logo_link, logo_alt, logo_height, logo_width, created_at, domain
`

type updateChangelogParams struct {
	Title       sql.NullString
	Subtitle    sql.NullString
	Subdomain   sql.NullString
	Domain      sql.NullString
	LogoSrc     sql.NullString
	LogoLink    sql.NullString
	LogoAlt     sql.NullString
	LogoHeight  sql.NullString
	LogoWidth   sql.NullString
	WorkspaceID string
	ID          string
}

func (q *Queries) updateChangelog(ctx context.Context, arg updateChangelogParams) (changelog, error) {
	row := q.db.QueryRowContext(ctx, updateChangelog,
		arg.Title,
		arg.Subtitle,
		arg.Subdomain,
		arg.Domain,
		arg.LogoSrc,
		arg.LogoLink,
		arg.LogoAlt,
		arg.LogoHeight,
		arg.LogoWidth,
		arg.WorkspaceID,
		arg.ID,
	)
	var i changelog
	err := row.Scan(
		&i.ID,
		&i.WorkspaceID,
		&i.Subdomain,
		&i.Title,
		&i.Subtitle,
		&i.SourceID,
		&i.LogoSrc,
		&i.LogoLink,
		&i.LogoAlt,
		&i.LogoHeight,
		&i.LogoWidth,
		&i.CreatedAt,
		&i.Domain,
	)
	return i, err
}
