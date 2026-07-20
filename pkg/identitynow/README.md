# SailPoint IdentityNow Adapter

## Entity External IDs

The Entity External ID is used directly as the API path segment. The adapter constructs URLs as:

```
{baseURL}/{apiVersion}/{entityExternalID}?limit={pageSize}&offset={offset}
```

The adapter is open-ended - any SailPoint IdentityNow API resource name can be used as an entity external ID.

| Entity External ID | API Path | Description |
|---|---|---|
| `accounts` | `/{apiVersion}/accounts` | SailPoint accounts |
| `entitlements` | `/{apiVersion}/entitlements` | Standalone entitlements |
| `accountEntitlements` | `/{apiVersion}/accounts/{accountID}/entitlements` | Entitlements scoped to a specific account (special case) |

## Special Cases

`accountEntitlements` is the only entity with custom URL construction. Instead of being appended directly as a path segment, it decomposes into a sub-resource pattern: `/accounts/{accountID}/entitlements`.

All other entity external IDs (e.g., `roles`, `identities`, `access-profiles`) are placed directly into the URL path.

## Supported API Versions

- `v3`
- `beta`

Each entity can optionally override the default API version via its entity configuration.

## Adding New Entity External IDs

**No code changes required for standard entities.** The adapter appends any entity external ID directly to the URL path. Simply configure a new entity in SGNL with the SailPoint resource name (e.g., `roles`, `identities`, `access-profiles`) and it will work.

To add a new entity with **custom URL construction** (like `accountEntitlements`):

1. Add a constant in `datasource.go`.
2. Add a case in `ConstructEndpoint` in `endpoint.go` with the custom URL pattern.
3. If the entity is a sub-resource (like account entitlements), implement collection-based pagination to iterate through the parent resource.
4. Add tests.
