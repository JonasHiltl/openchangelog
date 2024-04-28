// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: tokens.sql

package store

import (
	"context"
)

const createToken = `-- name: CreateToken :one
INSERT INTO tokens (
    key, workspace_id
) VALUES (
    $1, $2
) RETURNING key, workspace_id
`

type CreateTokenParams struct {
	Key         string
	WorkspaceID string
}

func (q *Queries) CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error) {
	row := q.db.QueryRow(ctx, createToken, arg.Key, arg.WorkspaceID)
	var i Token
	err := row.Scan(&i.Key, &i.WorkspaceID)
	return i, err
}

const getToken = `-- name: GetToken :one
SELECT key, workspace_id FROM tokens
WHERE key = $1
`

func (q *Queries) GetToken(ctx context.Context, key string) (Token, error) {
	row := q.db.QueryRow(ctx, getToken, key)
	var i Token
	err := row.Scan(&i.Key, &i.WorkspaceID)
	return i, err
}

const listToken = `-- name: ListToken :many
SELECT key, workspace_id FROM tokens
WHERE workspace_id = $1
`

func (q *Queries) ListToken(ctx context.Context, workspaceID string) ([]Token, error) {
	rows, err := q.db.Query(ctx, listToken, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Token
	for rows.Next() {
		var i Token
		if err := rows.Scan(&i.Key, &i.WorkspaceID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}