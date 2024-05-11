# Multi Tenancy

Openchangelog supports multi tenancy by storing `workspaces`, `sources` & `changelogs` in Postgres.  
**Note**: We do **not** store the changelog articles in Postgres, they still must be stored in a source like Github or locally.

To configure postgres and enable Multi Tenancy you need to specify the Postgres URL on the `config.yaml`.
```
# config.yaml
databaseUrl:
```

You can render the changelog of a specific workspace by calling `GET /:wid/:cid` and Openchangelog will fetch the config & source of the specified changelog.

To interact with `workspaces`, `sources` & `changelogs` you can use the REST API under the `/api/` endpoint.  
You need to authenticate on every request with a specific workspace by using the supplied `bearer` token when creating/returning your workspace.