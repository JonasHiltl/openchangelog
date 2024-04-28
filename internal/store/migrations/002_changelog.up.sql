CREATE TYPE source_type AS ENUM('GitHub');

CREATE TABLE IF NOT EXISTS changelogs (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    workspace_id VARCHAR(23) NOT NULL,
    title TEXT,
    subtitle TEXT,
    source_id BIGINT,
    source_type source_type,
    logo_src TEXT,
    logo_link TEXT,
    logo_alt TEXT,
    logo_height TEXT,
    logo_width TEXT,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (workspace_id, id)
);

ALTER TABLE changelogs
ADD CONSTRAINT fk_changelog_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE;