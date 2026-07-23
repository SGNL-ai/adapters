# Jira Cloud Adapter

## Entity External IDs

The Entity External ID maps to a Jira Cloud REST API endpoint. The adapter constructs URLs based on three base patterns:

- Default: `{baseURL}/rest/api/3/{endpoint}?...`
- Workspace: `{baseURL}/rest/servicedeskapi/assets/{endpoint}?...`
- Object: `{assetBaseURL}/workspace/{workspaceID}/v1/object/aql?...`

| Entity External ID | Unique ID Attribute | API Endpoint | Description |
|---|---|---|---|
| `User` | `accountId` | `/rest/api/3/users/search` | Jira users |
| `Issue` | `id` | `/rest/api/3/search` | Jira issues (JQL search) |
| `EnhancedIssue` | `id` | `/rest/api/3/search/jql` | Jira issues (new JQL endpoint with token pagination) |
| `Group` | `groupId` | `/rest/api/3/group/bulk` | Jira groups |
| `GroupMember` | synthetic `groupId-accountId` | `/rest/api/3/group/member?groupId={id}` | Group members (child of Group) |
| `Workspace` | `workspaceId` | `/rest/servicedeskapi/assets/workspace` | Jira Assets workspaces |
| `Object` | `globalId` | `{assetBaseURL}/workspace/{id}/v1/object/aql` | Jira Assets objects (child of Workspace, POST with AQL query) |

## Parent-Child Relationships

- `GroupMember` is a child of `Group` - iterates groups to fetch members per group.
- `Object` is a child of `Workspace` - iterates workspaces to fetch objects per workspace.

## Notes

- Issue queries use JQL (Jira Query Language) via the `filter` configuration.
- Object queries use AQL (Asset Query Language) via POST body.
- `EnhancedIssue` is a newer endpoint using token-based pagination instead of offset.

## Adding New Entity External IDs

**Requires code changes.** The adapter uses a fixed `ValidEntityExternalIDs` map with entity-specific endpoint paths, response parsing, and pagination strategies.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `Project string = "Project"`).
2. Add an entry to `ValidEntityExternalIDs` with the unique ID attribute, endpoint path, API base, and pagination type.
3. Add URL construction logic in `ConstructURL` for the new entity.
4. If the entity uses a different response structure, update the response parsing logic.
5. If it's a child entity, configure `MemberOf` to point to the parent and implement collection iteration.
6. Add tests.
