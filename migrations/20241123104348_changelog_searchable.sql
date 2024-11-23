-- +goose Up
-- +goose StatementBegin
ALTER TABLE changelogs ADD searchable INTEGER NOT NULL DEFAULT 0 check (searchable in (0, 1));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE changelogs DROP searchable;
-- +goose StatementEnd
