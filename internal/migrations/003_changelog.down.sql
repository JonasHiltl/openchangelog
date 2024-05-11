ALTER TABLE changelogs DROP CONSTRAINT fk_changelog_workspace;

DELETE TABLE changelogs;
DROP TYPE source_type;