# Workday Adapter

## Entity External IDs

The Entity External ID maps to a WQL (Workday Query Language) table name. The adapter constructs URLs as:

```
{baseURL}/api/wql/{apiVersion}/{organizationID}/data?limit={pageSize}&offset={offset}&query={WQL}
```

The WQL query uses the entity external ID in the `FROM` clause:

```sql
SELECT {columns} FROM {entityExternalID} [ORDER BY {uniqueIdColumn} ASC]
```

The adapter is open-ended - any valid WQL table name can be used as an entity external ID.

| Entity External ID | WQL Table | Description |
|---|---|---|
| `allWorkers` | allWorkers | All worker records |

## Attribute External IDs

Attributes use either plain column names or JSONPath expressions:

| Format | Example | WQL Column |
|---|---|---|
| Plain name | `employeeID` | `employeeID` |
| JSONPath | `$.worker.id` | `worker` (first path segment) |
| JSONPath | `$.managementLevel.descriptor` | `managementLevel` |
| JSONPath array | `$.email_Work[0].descriptor` | `email_Work` |

The `AttrExternalIDToColumnName` function converts JSONPath attributes to their top-level column name for the SELECT clause.

## Example

Entity External ID `allWorkers` with various attributes produces:

```
https://instance.workday.com/api/wql/v1/SGNL/data?limit=100&offset=0&query=SELECT+FTE,+company,+email_Work,+employeeID,+gender,+hireDate,+jobTitle,+managementLevel,+positionID,+worker,+workerActive+FROM+allWorkers
```

## Notes

- The only supported API version is `v1`.
- Max page size is 1000.
- Child entity external IDs (e.g., `email_Work`) are included as columns in the parent's SELECT query.
- Columns are sorted alphabetically in the SELECT clause.

## Adding New Entity External IDs

**No code changes required.** The adapter places any entity external ID directly into the WQL `FROM` clause. Simply configure a new entity in SGNL with the WQL table name (e.g., `allActiveEmployees`, `allCompanies`, `allLocations`) and it will work.

To find valid WQL table names, refer to Workday's WQL Data Dictionary or use the Workday REST API to list available data sources.
