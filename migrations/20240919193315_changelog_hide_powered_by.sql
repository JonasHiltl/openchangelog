-- +goose Up
-- +goose StatementBegin
ALTER TABLE changelogs ADD hide_powered_by INTEGER NOT NULL DEFAULT 0 check (hide_powered_by in (0, 1));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE changelogs DROP hide_powered_by;
-- +goose StatementEnd
