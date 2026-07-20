# ServiceNow Adapter

## Entity External IDs

The Entity External ID maps to a ServiceNow table name. The adapter constructs URLs as:

```
{baseURL}/api/now/{apiVersion}/table/{entityExternalID}?sysparm_fields={fields}&sysparm_limit={pageSize}
```

| Entity External ID | ServiceNow Table | Description |
|---|---|---|
| `sys_user` | sys_user | ServiceNow users |
| `sys_user_group` | sys_user_group | User groups |
| `sys_user_grmember` | sys_user_grmember | Group membership records |
| `sn_customerservice_case` | sn_customerservice_case | Customer service cases |
| `incident` | incident | Incidents |
| `change_request` | change_request | Change requests |
| `change_task` | change_task | Change tasks |

## How It Works

The Entity External ID is placed directly into the URL as the table name. ServiceNow's Table API returns all records from that table wrapped in a `"result"` JSON array.

Example URL:
```
https://instance.service-now.com/api/now/v2/table/change_request?sysparm_fields=sys_id,number,state&sysparm_limit=100&sysparm_query=ORDERBYsys_id
```

## Response Structure

ServiceNow always wraps its response in a top-level `result` key:

```json
{
  "result": [
    {"sys_id": "...", "number": "CHG0000001", ...},
    ...
  ]
}
```

The adapter parses objects from within `result` - attribute external IDs correspond to the JSON keys inside each object (e.g., `sys_id`, `state`, `assigned_to`).

## Notes

- Supports optional custom URL paths instead of the default `/api/now` prefix.
- Pagination uses ServiceNow's `Link` header with `rel="next"` for cursor-based pagination.
- The adapter is open-ended - any valid ServiceNow table name can be used as an entity external ID.
- All requests include `sysparm_exclude_reference_link=true` and order by `sys_id`.

## Adding New Entity External IDs

**No code changes required.** The adapter places any entity external ID directly into the URL path as the table name. Simply configure a new entity in SGNL with the ServiceNow table name (e.g., `problem`, `kb_knowledge`, `cmdb_ci_server`) and it will work.

The predefined constants in `datasource.go` (`sys_user`, `sys_user_group`, etc.) exist only for internal reference - the adapter does not validate against them. Any string passed as the entity external ID is used as-is in the API URL.
