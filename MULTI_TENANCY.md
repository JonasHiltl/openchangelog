# Multi Tenancy

Openchangelog supports multi tenancy by storing `workspaces`, `sources` & `changelogs` in SQLite.  
**Note**: We do **not** store the changelog articles in SQLite, they still must be stored in a source like Github.

To configure sqlite and enable Multi Tenancy you need to specify the SQLite URL on the `openchangelog.yml`.
Also set `?_foreign_keys=on` on the sqlite url to enfore foreign key constraints.
```
# openchangelog.yml
sqliteUrl:
```

You can render the changelog of a specific workspace by accessing it through the changelog's subdomain or host.

To interact with `workspaces`, `sources` & `changelogs` you can use the REST API under the `/api/` endpoint.  
You need to authenticate on every request with a specific workspace by using the supplied `bearer` token when creating/returning your workspace.