-- +goose Up
-- +goose StatementBegin
-- sqlc embed currently not workign with nullable embeds see this https://github.com/sqlc-dev/sqlc/issues/2997
-- this view is used, because sqlc treats views as nullable
CREATE VIEW changelog_source AS
SELECT gh.* -- in future probably need to prefix this with gh_
FROM changelogs cl
LEFT JOIN gh_sources gh 
    ON cl.workspace_id = gh.workspace_id
    AND cl.source_id LIKE 'gh_%'
    AND cl.source_id = gh.id
GROUP BY source_id, gh.workspace_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW changelog_source;
-- +goose StatementEnd
