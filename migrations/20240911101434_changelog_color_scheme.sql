-- +goose Up
-- +goose StatementBegin
ALTER TABLE changelogs ADD color_scheme INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE changelogs DROP color_scheme;
-- +goose StatementEnd
