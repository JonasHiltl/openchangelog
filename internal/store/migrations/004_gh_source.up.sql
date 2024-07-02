CREATE TABLE IF NOT EXISTS gh_sources (
    id TEXT NOT NULL,
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo TEXT NOT NULL,
    path TEXT NOT NULL,
    installation_id INTEGER NOT NULL,
    PRIMARY KEY (workspace_id, id)
);