# HashiCorp Boundary Adapter

## Entity External IDs

The Entity External ID is used directly as the API path segment. The adapter constructs URLs as:

```
{baseURL}/v1/{entityExternalID}?page_size={n}&recursive=true&scope_id={scope}
```

| Entity External ID | Parent Entity | Description |
|---|---|---|
| `hosts` | `host-catalogs` | Boundary hosts (fetched per host catalog) |
| `host-sets` | `host-catalogs` | Boundary host sets (fetched per host catalog) |
| `credentials` | `credential-stores` | Stored credentials (fetched per credential store) |
| `credential-libraries` | `credential-stores` | Credential libraries (fetched per credential store) |
| `accounts` | `auth-methods` | User accounts (fetched per auth method) |
| `host-catalogs` | (none) | Host catalogs (top-level) |
| `credential-stores` | (none) | Credential stores (top-level) |
| `auth-methods` | (none) | Authentication methods (top-level) |

## Parent-Child Relationships

Five entities are children that require iterating through a parent collection first:

| Child | Parent | Query Param Added |
|---|---|---|
| `hosts` | `host-catalogs` | `host_catalog_id={parentID}` |
| `host-sets` | `host-catalogs` | `host_catalog_id={parentID}` |
| `credentials` | `credential-stores` | `credential_store_id={parentID}` |
| `credential-libraries` | `credential-stores` | `credential_store_id={parentID}` |
| `accounts` | `auth-methods` | `auth_method_id={parentID}` |

## Notes

- The adapter is open-ended - any valid Boundary resource type (e.g., `roles`, `users`, `groups`) can be used as an entity external ID.
- Top-level entities are fetched directly without parent iteration.
- Pagination uses `list_token` for cursor-based pagination.

## Adding New Entity External IDs

**No code changes required for top-level entities.** The adapter appends any entity external ID directly to the URL path (`/v1/{entityExternalID}`). Simply configure a new entity in SGNL with the Boundary resource type name (e.g., `roles`, `users`, `groups`, `scopes`) and it will work.

To add a new **child** entity (one that requires parent iteration):

1. Add a constant in `datasource.go` for the child entity.
2. Add a routing case in the `GetPage` function that maps the child to its parent entity and the query parameter to include (following the pattern of `hosts` -> `host-catalogs` with `host_catalog_id`).
3. Add the entity to `CursorsWithParentEntity` in `validation.go`.
4. Add tests.
