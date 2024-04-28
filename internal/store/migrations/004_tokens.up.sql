CREATE TABLE IF NOT EXISTS tokens (
    key VARCHAR(20) PRIMARY KEY,
    workspace_id VARCHAR(23) NOT NULL
);

ALTER TABLE tokens
ADD CONSTRAINT fk_tokens_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE;