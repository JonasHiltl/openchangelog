-- stored source does not support access token authentication
CREATE TABLE IF NOT EXISTS gh_sources (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    workspace_id VARCHAR(23) NOT NULL,
    owner VARCHAR(39) NOT NULL,
    repo VARCHAR(200) NOT NULL,
    path TEXT NOT NULL,
    installation_id BIGINT NOT NULL,
    PRIMARY KEY (workspace_id, id)
);

ALTER TABLE gh_sources
ADD CONSTRAINT fk_gh_sources_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE;