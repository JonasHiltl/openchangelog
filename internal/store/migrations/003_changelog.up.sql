CREATE TABLE IF NOT EXISTS changelogs (
    id TEXT NOT NULL,
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
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