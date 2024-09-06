-- +goose Up
-- +goose StatementBegin
ALTER TABLE changelogs ADD domain TEXT;
CREATE UNIQUE INDEX changelogs_domain ON changelogs(domain);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX changelogs_domain;
ALTER TABLE changelogs DROP domain;
-- +goose StatementEnd
