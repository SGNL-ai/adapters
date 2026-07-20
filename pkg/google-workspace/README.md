# Google Workspace Adapter

## Entity External IDs

The Entity External ID maps to a Google Admin Directory API path. The adapter constructs URLs as:

```
{baseURL}/admin/directory/{apiVersion}/{path}?maxResults={pageSize}&domain={domain}
```

| Entity External ID | API Path | Unique ID Attribute | Max Page Size | Description |
|---|---|---|---|---|
| `User` | `/admin/directory/v1/users` | `id` | 500 | Google Workspace users |
| `Group` | `/admin/directory/v1/groups` | `id` | 1000 | Google Workspace groups |
| `Member` | `/admin/directory/v1/groups/{groupID}/members` | `uniqueId` | 1000 | Group members (child of Group) |

## Parent-Child Relationships

- `Member` is a child entity - the adapter iterates through groups and fetches members per group.
- Member's `uniqueId` is synthesized as `"{groupId}-{memberId}"`.

## Response Parsing

Each entity's response is parsed from a different JSON field:
- `User` → `response.Users`
- `Group` → `response.Groups`
- `Member` → `response.Members`

## Notes

- Supports `domain` or `customer` scoping via query parameters.
- User and Group support ordering by `EMAIL`.
- Member requires `id` and `groupId` attributes in the request.

## Adding New Entity External IDs

**Requires code changes.** The adapter validates against `ValidEntityExternalIDs` and uses entity-specific URL construction and response parsing.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `OrgUnit = "OrgUnit"`).
2. Add an entry to `ValidEntityExternalIDs` with the path template, unique ID attribute, max page size, order-by attribute, and any required attributes.
3. Add a case in `ConstructEndpoint` in `endpoint.go` for entity-specific query parameters (similar to `AddUserParams`/`AddGroupParams`).
4. Update `ParseResponse` in `datasource.go` to extract objects from the correct JSON field in the Google API response.
5. If the entity is a child (like Member), configure the `MemberOf` relationship and collection-based pagination.
6. Add tests.
