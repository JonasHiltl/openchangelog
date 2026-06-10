-- +goose Up
ALTER TABLE changelogs ADD COLUMN hide_rss_icon INTEGER NOT NULL DEFAULT 0;
-- +goose Down
ALTER TABLE changelogs DROP COLUMN hide_rss_icon;
