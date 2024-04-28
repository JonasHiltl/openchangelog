CREATE VIEW changelog_sources AS (
  SELECT gh.* FROM changelogs cl
  LEFT JOIN gh_sources gh ON cl.source_type = 'GitHub' AND gh.id = cl.source_id
);