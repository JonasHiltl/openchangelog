-- +goose Up
-- +goose StatementBegin
ALTER TABLE changelogs ADD analytics INTEGER NOT NULL DEFAULT 0 check (analytics in (0, 1));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE changelogs DROP analytics;
-- +goose StatementEnd
