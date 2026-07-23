# DB2 Adapter

## Entity External IDs

The Entity External ID maps to a DB2 table name. The adapter constructs SQL queries as:

```sql
SELECT {columns} FROM [{schema}.]{entityExternalID}
WHERE {cursor_condition}
ORDER BY {unique_key_columns}
FETCH FIRST {pageSize + 1} ROWS ONLY
```

The adapter is open-ended - any valid DB2 table name can be used as an entity external ID.

| Entity External ID | SQL Table | Description |
|---|---|---|
| `ITEMS` | `"ITEMS"` | Items table |
| `EKPO` | `"EKPO"` | SAP purchase order items |
| `users` | `"users"` | Users table |

## How It Works

The Entity External ID is quoted and placed directly in the SQL `FROM` clause. If a schema is configured, it becomes `"SCHEMA"."TABLE"`.

Example with entity external ID `ITEMS`, schema `TEST_SCHEMA`:
```sql
SELECT "TENANT_ID", "DOC_NUM", "AMOUNT", "REGION"
FROM "TEST_SCHEMA"."ITEMS"
WHERE CAST("TENANT_ID" AS VARCHAR(50)) > ?
ORDER BY CAST("TENANT_ID" AS VARCHAR(50))
FETCH FIRST 11 ROWS ONLY
```

## Synthetic Composite ID

When an entity has the attribute external ID `id` marked as the unique ID, the adapter:
1. Queries DB2 system catalog (`SYSCAT.TABCONST` + `SYSCAT.KEYCOLUSE`) to discover primary key columns.
2. Adds those PK columns to the SELECT.
3. Concatenates their values with `|` as separator (e.g., `"T1|D1003|L01"`).

## Connection

The adapter connects via TCP using the DB2 driver:
```
HOSTNAME={host};DATABASE={db};UID={user};PWD={pass};PORT={port};PROTOCOL=TCPIP
```

## Notes

- All entity and attribute external IDs are validated as safe SQL identifiers before use.
- Pagination uses cursor-based keyset pagination (comparing cast-to-VARCHAR values).
- Filters are configured per entity external ID in the adapter configuration.
- Supports optional SSL connections.

## Adding New Entity External IDs

**No code changes required.** The adapter uses any entity external ID directly as the DB2 table name in the `FROM` clause. Simply configure a new entity in SGNL with the DB2 table name (e.g., `EMPLOYEES`, `ORDERS`, `MARA`) and it will work.

Requirements for the entity external ID:
- Must be a valid SQL identifier (letters, digits, `$`, `_`, max 128 characters).
- Must correspond to an existing table in the configured database (and schema, if specified).
- Attribute external IDs must match column names in that table.
