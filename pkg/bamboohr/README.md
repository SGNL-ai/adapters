# BambooHR Adapter

## Entity External IDs

The Entity External ID maps to a BambooHR API path. The adapter constructs URLs as:

```
{baseURL}/{companyDomain}/{apiVersion}/reports/custom?format=JSON&onlyCurrent={bool}
```

| Entity External ID | API Path | Description |
|---|---|---|
| `Employee` | `/reports/custom` | BambooHR employee records |

## How It Works

This is the only supported entity. The adapter issues a POST request with a JSON body containing the list of attribute external IDs (field names) to retrieve:

```json
{"fields": ["id", "firstName", "lastName", ...]}
```

## Required Attribute

Every request must include the attribute with external ID `id` - this is the unique identifier for employee records.

## Adding New Entity External IDs

**Requires code changes.** The adapter only supports `Employee` via the `ValidEntityExternalIDs` map, and the URL construction is hardcoded to the Custom Reports API path.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `TimeOff = "TimeOff"`).
2. Add an entry to `ValidEntityExternalIDs` with the corresponding BambooHR API path.
3. Update `ConstructEndpoint` in `endpoint.go` to handle the new path pattern.
4. If the new entity uses a different request method or body format, update `GetPage` in `datasource.go`.
5. Add validation for the new entity's required unique ID attribute in `validation.go`.
6. Add tests.
