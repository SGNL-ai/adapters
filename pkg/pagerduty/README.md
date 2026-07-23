# PagerDuty Adapter

## Entity External IDs

The Entity External ID maps to a PagerDuty API path segment. The adapter constructs URLs as:

```
https://api.pagerduty.com/{entityExternalID}?offset={cursor}&limit={pageSize}
```

| Entity External ID | API Endpoint | Unique ID Attribute | Description |
|---|---|---|---|
| `users` | `GET /users` | `id` | PagerDuty users |
| `teams` | `GET /teams` | `id` | PagerDuty teams |
| `oncalls` | `GET /oncalls` | composite | On-call schedules |
| `members` | `GET /teams/{teamID}/members` | composite | Team members (child of teams) |

## Parent-Child Relationships

- `members` is a child of `teams` - the adapter iterates through teams and fetches members per team.

## Composite Unique IDs

- `members`: Synthesized as `{teamID}-{userID}` since a user can belong to multiple teams.
- `oncalls`: Synthesized as `{escalation_policy.id}-{user.id}-{start}-{end}` since on-call objects lack a native unique ID.

## Notes

- The base URL is always `https://api.pagerduty.com` (enforced by validation).
- Pagination uses offset-based pagination.
- Additional query parameters can be configured per entity (e.g., `include[]=contact_methods&include[]=teams` for users).
- Max page size is 100.

## Adding New Entity External IDs

**Requires code changes for child entities; top-level entities work without changes.** The adapter appends the entity external ID directly as a URL path segment for non-member entities.

For a new **top-level** entity (like `services`, `escalation_policies`):
- No code changes needed. Configure the entity in SGNL with the PagerDuty resource name and it will be fetched from `https://api.pagerduty.com/{entityExternalID}`.

For a new **child** entity (scoped to a parent, like `members` is scoped to `teams`):

1. Add routing logic in `datasource.go` to construct the parent-scoped URL (e.g., `/teams/{parentID}/members`).
2. Implement collection-based pagination to iterate through the parent entity.
3. Implement composite unique ID synthesis (since child objects typically lack a globally unique ID).
4. Add tests.
