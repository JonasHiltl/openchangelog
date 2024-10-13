-- +goose Up
-- +goose StatementBegin
ALTER TABLE changelogs ADD protected INTEGER NOT NULL DEFAULT 0 check (protected in (0, 1));
ALTER TABLE changelogs ADD password_hash TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE changelogs DROP protected;
ALTER TABLE changelogs DROP password_hash;
-- +goose StatementEnd
