-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS changelogs (
    id TEXT NOT NULL,
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    subdomain TEXT NOT NULL,
    title TEXT,
    subtitle TEXT,
    source_id TEXT,
    logo_src TEXT,
    logo_link TEXT,
    logo_alt TEXT,
    logo_height TEXT,
    logo_width TEXT,
    created_at INTEGER NOT NULL DEFAULT (unixepoch('now')),
    PRIMARY KEY (workspace_id, id)
) STRICT;

CREATE UNIQUE INDEX changelogs_subdomain ON changelogs(subdomain);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX changelogs_subdomain;
DROP TABLE changelogs;
-- +goose StatementEnd
