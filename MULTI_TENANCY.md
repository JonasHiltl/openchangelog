# Multi Tenancy

Openchangelog supports multi tenancy by storing `workspaces`, `sources` & `changelogs` in SQLite.  
**Note**: We do **not** store the changelog articles in SQLite, they still must be stored in a source like Github.

To configure sqlite and enable Multi Tenancy you need to specify the SQLite URL on the `openchangelog.yml`.
```
# openchangelog.yml
sqliteUrl:
```

You can render the changelog of a specific workspace by calling `GET /?wid=...&cid=...` and Openchangelog will fetch the changelog & source of the specified changelog from SQLite.

To interact with `workspaces`, `sources` & `changelogs` you can use the REST API under the `/api/` endpoint.  
You need to authenticate on every request with a specific workspace by using the supplied `bearer` token when creating/returning your workspace.