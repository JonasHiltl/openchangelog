-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workspaces (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
) STRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workspaces;
-- +goose StatementEnd
