# Okta Adapter

## Entity External IDs

The Entity External ID maps to an Okta API resource path. The adapter constructs URLs as:

```
{baseURL}/api/v1/{resource}?limit={pageSize}
```

| Entity External ID | API Path | Unique ID Attribute | Description |
|---|---|---|---|
| `User` | `/api/v1/users` | `id` | Okta users |
| `Group` | `/api/v1/groups` | `id` | Okta groups |
| `GroupMember` | `/api/v1/groups/{groupID}/users` | `id` | Group members (child of Group) |
| `Application` | `/api/v1/apps` | `id` | Okta applications |

## Parent-Child Relationships

- `GroupMember` is a child of `Group` - the adapter iterates through groups and fetches members per group.

## Filter and Search Support

| Entity | Filter | Search | Default Filter |
|---|---|---|---|
| `User` | Yes | Yes (mutually exclusive with filter) | None |
| `Group` | Yes | Yes (mutually exclusive with filter) | `type eq "OKTA_GROUP" or type eq "APP_GROUP"` |
| `GroupMember` | No | No | None |
| `Application` | Yes | No | None |

## Notes

- Group has a default filter applied when neither filter nor search is configured.
- Pagination uses Link header-based cursors (Okta's standard pagination).

## Adding New Entity External IDs

**Requires code changes.** The adapter validates against `ValidEntityExternalIDs` and uses a switch statement in `endpoint.go` to construct entity-specific URL paths.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `Role string = "Role"`).
2. Add the entity to `ValidEntityExternalIDs`.
3. Add a case in the `switch` statement in `ConstructEndpoint` in `endpoint.go` with the Okta API path.
4. If the entity supports filter or search, add handling in the endpoint construction.
5. If it's a child entity (like GroupMember), implement collection-based pagination with a parent reference.
6. Add tests.
