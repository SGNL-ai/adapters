# Rootly Adapter

## Entity External IDs

The Entity External ID is used directly as the API path segment. The adapter constructs URLs as:

```
{address}/{apiVersion}/{entityExternalID}?page[size]={pageSize}&page[number]={page}
```

The adapter is open-ended - any Rootly API resource name can be used as an entity external ID.

| Entity External ID | API Endpoint | Description |
|---|---|---|
| `incidents` | `/v1/incidents` | Rootly incidents |
| `users` | `/v1/users` | Rootly users |
| `teams` | `/v1/teams` | Rootly teams |

## How It Works

The Entity External ID is appended directly as a path segment after the base URL and API version. Any valid Rootly API resource name will work.

Example URL:
```
https://api.rootly.com/v1/incidents?page[size]=100&page[number]=1&filter[status]=started&include=form_field_values
```

## Configuration

Filters and includes are configured per entity using the entity external ID as the map key:

```json
{
  "apiVersion": "v1",
  "filters": {
    "incidents": "status=started&severity=high",
    "users": "email=user@example.com"
  },
  "includes": {
    "incidents": "roles",
    "users": "role,email_addresses"
  }
}
```

## Notes

- The only supported API version is `v1`.
- Pagination uses page-number-based pagination (`page[size]` + `page[number]`).

## Adding New Entity External IDs

**No code changes required.** The adapter appends any entity external ID directly as a URL path segment. Simply configure a new entity in SGNL with the Rootly API resource name (e.g., `services`, `environments`, `functionalities`, `severities`) and it will work.

Filters and includes for the new entity can be added via the adapter configuration using the entity external ID as the map key.
