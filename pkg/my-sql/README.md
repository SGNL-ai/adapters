# MySQL Adapter

## Entity External IDs

The Entity External ID maps to a MySQL table name. The adapter constructs SQL queries as:

```sql
SELECT *, CAST({uniqueAttribute} AS CHAR(50)) AS str_id, COUNT(*) OVER() AS total_remaining_rows
FROM {entityExternalID}
WHERE CAST({uniqueAttribute} AS CHAR(50)) > ?
ORDER BY str_id ASC
LIMIT ?
```

The adapter is open-ended - any valid MySQL table name can be used as an entity external ID.

| Entity External ID | SQL Table | Description |
|---|---|---|
| `users` | `users` | Users table |
| `groups` | `groups` | Groups table |

## How It Works

The Entity External ID is used directly as the table name in the SQL `FROM` clause. Attribute external IDs map to column names in the result set.

Example with entity external ID `users`, unique attribute `id`:
```sql
SELECT *, CAST(`id` AS CHAR(50)) AS `str_id`, COUNT(*) OVER() AS `total_remaining_rows`
FROM `users`
WHERE (CAST(`id` AS CHAR(50)) > ?)
ORDER BY `str_id` ASC
LIMIT ?
```

## Connection

The adapter connects via TCP using the MySQL driver:
```
{username}:{password}@tcp({host:port})/{database}
```

## Validation

Entity external IDs (table names) and unique attribute external IDs are validated against:
```
^[a-zA-Z0-9$_]{1,128}$
```

This prevents SQL injection since table names cannot be parameterized in prepared statements.

## Notes

- Pagination uses cursor-based keyset pagination (comparing cast-to-CHAR values).
- Filters are configured per entity external ID in the adapter configuration.
- The adapter uses the `goqu` query builder library for SQL construction.
- Two versions exist: `0.0.1-alpha` and `0.0.2-alpha`.

## Adding New Entity External IDs

**No code changes required.** The adapter uses any entity external ID directly as the MySQL table name in the `FROM` clause. Simply configure a new entity in SGNL with the MySQL table name (e.g., `employees`, `orders`, `departments`) and it will work.

Requirements for the entity external ID:
- Must match the regex `^[a-zA-Z0-9$_]{1,128}$` (validated to prevent SQL injection).
- Must correspond to an existing table in the configured database.
- Attribute external IDs must match column names in that table.
