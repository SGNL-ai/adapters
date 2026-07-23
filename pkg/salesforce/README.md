# Salesforce Adapter

## Entity External IDs

The Entity External ID maps to a Salesforce object (table) name. The adapter constructs SOQL queries as:

```
SELECT Id,{attributes} FROM {entityExternalID} [WHERE {filter}] ORDER BY Id ASC
```

Delivered via the Salesforce REST API:
```
{baseURL}/services/data/v{apiVersion}/query?q={SOQL}
```

The adapter is open-ended - any valid Salesforce object name (standard or custom) can be used as an entity external ID.

| Entity External ID | Salesforce Object | Description |
|---|---|---|
| `Account` | Account | Salesforce accounts |
| `Contact` | Contact | Salesforce contacts |
| `User` | User | Salesforce users |
| `Case` | Case | Support cases |
| `CustomObject__c` | (any custom object) | Custom Salesforce objects |

## How It Works

The Entity External ID is URL-escaped and placed directly into the SOQL `FROM` clause. The adapter builds a `SELECT` statement from the configured attribute external IDs plus the required `Id` field.

Example URL (first page):
```
https://instance.salesforce.com/services/data/v58.0/query?q=SELECT+Id,CaseNumber,Status+FROM+Case+WHERE+Status+%3D+%27Open%27+ORDER+BY+Id+ASC
```

Subsequent pages use the cursor URL returned by Salesforce:
```
https://instance.salesforce.com/services/data/v58.0/query/0r8Hu1lKCluUiC9IMK-200
```

## Required Attribute

Every request must include the attribute with external ID `Id` - this is Salesforce's universal unique identifier.

## Attribute External IDs

Attribute external IDs can be:
- Plain field names: `Id`, `Name`, `CustomField__c`
- JSONPath format for relationship traversal: `$.Account.Name` (up to 5 levels)

## Supported API Versions

`52.0` through `58.0`

## Notes

- Page size range: 200-2000.
- Filters are keyed by entity external ID in the adapter configuration.

## Adding New Entity External IDs

**No code changes required.** The adapter places any entity external ID directly into the SOQL `FROM` clause. Simply configure a new entity in SGNL with the Salesforce object API name (e.g., `Opportunity`, `Lead`, `Task`, or a custom object like `MyObject__c`) and it will work.

To find valid object names, check the Salesforce Object Reference or use the Describe API (`/services/data/v{version}/sobjects/`). Custom objects always end with `__c`.
