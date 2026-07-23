# CrowdStrike Adapter

## Entity External IDs

The CrowdStrike adapter supports two API mechanisms: GraphQL (Identity Protection) and REST (Endpoint Protection). The Entity External ID determines which mechanism and endpoint is used.

### GraphQL Entities (Identity Protection)

These entities query the CrowdStrike Identity Protection GraphQL API:

```
{baseURL}/identity-protection/combined/graphql/{apiVersion}
```

| Entity External ID | Unique ID Attribute | Order By | Description |
|---|---|---|---|
| `user` | `entityId` | `RISK_SCORE` | Identity Protection users |
| `incident` | `incidentId` | `END_TIME` | Identity Protection incidents |
| `endpoint` | `entityId` | `RISK_SCORE` | Identity Protection endpoints |

### How GraphQL Queries Work

Each entity external ID has a dedicated `QueryBuilder` that generates a specific GraphQL query. The adapter uses the `machinebox/graphql` client to POST queries to the Identity Protection endpoint.

**Query dispatch** (`GetQueryBuilder` in `query.go`):
- `user` -> `UserQueryBuilder`
- `incident` -> `IncidentQueryBuilder`
- `endpoint` -> `EndpointQueryBuilder`

#### User Query

Queries the `entities` field with `types: [USER]`:

```graphql
{
  entities(
    archived: false
    enabled: true
    types: [USER]
    sortKey: RISK_SCORE
    sortOrder: DESCENDING
    first: 100
    after: "cursor..."
  ) {
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      ... on UserEntity {
        entityId
        primaryDisplayName
        emailAddresses
        riskScore
        riskScoreSeverity
        mostRecentActivity
        accounts {
          ... on ActiveDirectoryAccountDescriptor {
            domain
            samAccountName
            upn
            enabled
            ...
          }
        }
        riskFactors {
          score
          severity
          type
        }
      }
    }
  }
}
```

Key fields returned: `entityId`, `primaryDisplayName`, `emailAddresses`, `riskScore`, `riskScoreSeverity`, `impactScore`, `mostRecentActivity`, `inactive`, `stale`, `watched`, `shared`, nested `accounts` (Active Directory details), and `riskFactors`.

#### Incident Query

Queries the `incidents` field:

```graphql
{
  incidents(
    first: 100
    sortKey: END_TIME
    sortOrder: DESCENDING
    after: "cursor..."
  ) {
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      incidentId
      severity
      type
      startTime
      endTime
      lifeCycleStage
      compromisedEntities { entityId, primaryDisplayName, riskScore, ... }
      alertEvents { alertId, alertType, eventSeverity, entities {...}, ... }
    }
  }
}
```

Key fields returned: `incidentId`, `severity`, `type`, `startTime`, `endTime`, `lifeCycleStage`, nested `compromisedEntities` (entities involved), and `alertEvents` (alerts with their own nested entities).

#### Endpoint Query

Queries the `entities` field with `types: [ENDPOINT]`:

```graphql
{
  entities(
    archived: false
    enabled: true
    types: [ENDPOINT]
    sortKey: RISK_SCORE
    sortOrder: DESCENDING
    first: 100
    after: "cursor..."
  ) {
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      ... on EndpointEntity {
        entityId
        hostName
        lastIpAddress
        riskScore
        ztaScore
        agentId
        accounts {
          ... on ActiveDirectoryAccountDescriptor { ... }
        }
        riskFactors { score, severity, type }
      }
    }
  }
}
```

Key fields returned: `entityId`, `hostName`, `lastIpAddress`, `riskScore`, `ztaScore`, `agentId`, `agentVersion`, `cid`, `unmanaged`, `staticIpAddresses`, nested `accounts` and `riskFactors`.

### GraphQL Pagination

All GraphQL entities use cursor-based pagination via the `pageInfo` object:
- `hasNextPage` - indicates whether more results exist
- `endCursor` - opaque cursor string passed as the `after` parameter in the next request

The cursor is base64-encoded and stored as a `CompositeCursor` between pages.

### GraphQL Response Parsing

The response is parsed based on the entity external ID:
- `user` and `endpoint` -> reads from `data.entities.nodes`
- `incident` -> reads from `data.incidents.nodes`

### Configuration Parameters

GraphQL entities support these config options:
- `archived` (bool) - include archived entities (applies to `user` and `endpoint`)
- `enabled` (bool) - filter to enabled entities only (applies to `user` and `endpoint`)
- `apiVersion` - API version for the endpoint URL path

### REST Entities (Endpoint Protection)

These entities query CrowdStrike REST APIs:

| Entity External ID | List Endpoint | Get Endpoint | Description |
|---|---|---|---|
| `endpoint_protection_device` | `GET devices/queries/devices-scroll/v1` | `POST devices/entities/devices/v2` | Falcon endpoint devices |
| `endpoint_protection_incident` | `GET incidents/queries/incidents/v1` | `POST incidents/entities/incidents/GET/v1` | Falcon incidents |
| `endpoint_protection_alert` | (combined API) | `POST alerts/combined/alerts/v1` | Falcon alerts |

## How REST Entities Work

- `endpoint_protection_device` and `endpoint_protection_incident` use a two-step pattern: first list resource IDs (GET), then fetch full objects by ID (POST with `{"ids": [...]}`).
- `endpoint_protection_alert` uses a single combined API with pagination via `after` cursor in the POST body.

## Adding New Entity External IDs

**Requires code changes.** Both GraphQL and REST entities are hardcoded with specific query builders and endpoint mappings.

### Adding a GraphQL entity:

1. Add a constant in `datasource.go` (e.g., `Group string = "group"`).
2. Add the entity to `ValidGraphQLEntityExternalIDs` in `datasource.go` with its unique ID attribute and order-by field.
3. Create a new `QueryBuilder` struct and `Build` method in `query.go` that generates the GraphQL query string.
4. Add a case in `GetQueryBuilder` to dispatch to the new builder.
5. If the response field differs from `entities`, update `parseGraphQLResponse` in `datasource_graphql.go`.
6. Add tests.

### Adding a REST entity:

1. Add a constant in `datasource.go`.
2. Add the entity to `ValidRESTEntityExternalIDs` in `datasource.go` specifying cursor type and whether it uses the list-then-get pattern.
3. Add an `EndpointInfo` entry in `EntityExternalIDToEndpoint` in `endpoint.go` with the list and/or get endpoint paths.
4. If the entity needs special request body handling (like alerts), add logic in `datasource_rest.go`.
5. Add tests.
