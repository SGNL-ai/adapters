# Duo Adapter

## Entity External IDs

The Entity External ID maps to a Duo Admin API path segment. The adapter constructs URLs as:

```
{baseURL}/admin/v1/{path}?limit={pageSize}&offset={offset}
```

| Entity External ID | API Path | Unique ID Attribute | Description |
|---|---|---|---|
| `User` | `users` | `user_id` | Duo users |
| `Group` | `groups` | `group_id` | Duo groups |
| `Phone` | `phones` | `phone_id` | Registered phones |
| `Endpoint` | `endpoints` | `epkey` | Registered endpoints/devices |

## Notes

- The only supported API version is `v1`.
- Pagination uses offset-based pagination (`limit` + `offset` query params).

## Adding New Entity External IDs

**Requires code changes.** The adapter validates against a fixed `ValidEntityExternalIDs` map that defines the path segment and unique ID attribute for each entity.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `Token = "Token"`).
2. Add an entry to `ValidEntityExternalIDs` with the Duo Admin API path (e.g., `tokens`) and unique ID attribute.
3. Validation will automatically accept the new entity since it checks against this map.
4. Add tests in `endpoint_test.go` and `adapter_test.go`.
