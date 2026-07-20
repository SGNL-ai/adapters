# SCIM 2.0 Adapter

## Entity External IDs

The Entity External ID is used directly as the API path segment. The adapter constructs URLs as:

```
{baseURL}/{entityExternalID}?startIndex={cursor}&count={pageSize}
```

The adapter is open-ended - any SCIM 2.0 resource type name can be used as an entity external ID.

| Entity External ID | API Endpoint | Description |
|---|---|---|
| `Users` | `{baseURL}/Users` | SCIM users |
| `Groups` | `{baseURL}/Groups` | SCIM groups |

## How It Works

The Entity External ID is appended as a path segment to the base URL. This follows the SCIM 2.0 standard where resource types are accessed by name in the URL path.

Example URL:
```
https://scim.example.com/Users?startIndex=1&count=100&filter=userType+eq+"Employee"&sortBy=userName&sortOrder=ascending
```

## Configuration

Query parameters can be configured per entity using the entity external ID as the map key:

```json
{
  "queryParams": {
    "Users": {
      "filter": "userType eq \"Employee\"",
      "sortBy": "userName",
      "ascending": true
    },
    "Groups": {
      "filter": "displayName co \"Engineering\"",
      "sortBy": "displayName",
      "ascending": true
    }
  }
}
```

## Notes

- Pagination uses SCIM's standard `startIndex`/`count` parameters (1-based indexing).
- Supports optional `filter`, `sortBy`, and `sortOrder` query parameters per entity.
- Any SCIM-compliant identity provider (Okta SCIM, Azure AD SCIM, etc.) can be used as the base URL.

## Adding New Entity External IDs

**No code changes required.** The adapter appends any entity external ID directly as a URL path segment. Simply configure a new entity in SGNL with the SCIM resource type name (e.g., `Users`, `Groups`, `Roles`, `Entitlements`) and it will work.

Query parameters for the new entity can be added via the adapter configuration using the entity external ID as the `queryParams` map key.
