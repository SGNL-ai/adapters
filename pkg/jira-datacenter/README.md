# Jira Data Center Adapter

## Entity External IDs

The Entity External ID maps to a Jira Data Center REST API endpoint. The adapter constructs URLs as:

```
{baseURL}/rest/api/{apiVersion}/{endpoint}?...
```

Where `apiVersion` defaults to `latest`.

| Entity External ID | Unique ID Attribute | API Endpoint | Response Array Field | Description |
|---|---|---|---|---|
| `User` | `key` | `/rest/api/{version}/group/member?groupname={groupID}` | `values` | Jira users (fetched via group membership, child of Group) |
| `Issue` | `id` | `/rest/api/{version}/search?jql={filter}` | `issues` | Jira issues |
| `Group` | `name` | `/rest/api/{version}/groups/picker` | `groups` | Jira groups |
| `GroupMember` | synthetic `groupID-userKey` | `/rest/api/{version}/group/member?groupname={groupID}` | `values` | Group members (child of Group) |

## Parent-Child Relationships

- `User` is a child of `Group` - users are fetched through group membership.
- `GroupMember` is a child of `Group` - members are fetched per group.

## Key Differences from Jira Cloud

- Uses configurable API version (`latest` or `2`) instead of fixed `v3`.
- `User` unique ID is `key` (not `accountId`).
- `Group` unique ID is `name` (not `groupId`).
- No Assets (Workspace/Object) support.
- No `EnhancedIssue` entity.
- Issue queries dynamically build a `fields` parameter from requested attributes.

## Adding New Entity External IDs

**Requires code changes.** The adapter validates against `ValidEntityExternalIDs` and uses entity-specific endpoint paths and response parsing.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `ProjectExternalID string = "Project"`).
2. Add an entry to `ValidEntityExternalIDs` with the unique ID attribute, endpoint path, response array field, and pagination config.
3. Add URL construction logic in the entity's `ConstructURL` method.
4. If the response array uses a different JSON field name, specify it in the entity config.
5. If it's a child entity, configure the parent relationship and collection iteration.
6. Add tests.
