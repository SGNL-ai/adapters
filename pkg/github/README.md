# GitHub Adapter

## Entity External IDs

The GitHub adapter primarily uses GraphQL to query GitHub Enterprise data, with one REST-based entity. The Entity External ID determines which GraphQL query or REST endpoint is used.

### GraphQL Entities

All GraphQL entities query:
- **Enterprise Cloud:** `{baseURL}/graphql`
- **Enterprise Server:** `{baseURL}/api/graphql`

| Entity External ID | Unique ID Attribute | Description |
|---|---|---|
| `Organization` | `id` | Enterprise organizations |
| `OrganizationUser` | `uniqueId` | Users within a specific organization (child of Organization) |
| `User` | `id` | All users across enterprise organizations |
| `Team` | `id` | Teams across organizations |
| `$.members.edges` (TeamMember) | `$.node.id` | Team members (child entity) |
| `$.repositories.edges` (TeamRepository) | `$.node.id` | Team repositories (child entity) |
| `Repository` | `id` | Repositories across organizations |
| `$.collaborators.edges` (RepositoryCollaborator) | `id` | Repository collaborators (child entity) |
| `Collaborator` | `id` | Collaborators with expanded traversal |
| `Label` | `id` | Repository labels |
| `IssueLabel` | `uniqueId` | Labels on issues (requires `labelId`, `id`) |
| `PullRequestLabel` | `uniqueId` | Labels on pull requests (requires `labelId`, `id`) |
| `Issue` | `id` | Repository issues |
| `IssueAssignee` | `uniqueId` | Issue assignees (requires `issueId`, `id`) |
| `IssueParticipant` | `uniqueId` | Issue participants (requires `issueId`, `id`) |
| `PullRequest` | `id` | Repository pull requests |
| `PullRequestChangedFile` | `uniqueId` | PR changed files (requires `pullRequestId`, `path`) |
| `PullRequestReview` | `id` | PR reviews |
| `PullRequestCommit` | `id` | PR commits |
| `PullRequestAssignee` | `uniqueId` | PR assignees (requires `pullRequestId`, `id`) |
| `PullRequestParticipant` | `uniqueId` | PR participants (requires `pullRequestId`, `id`) |

### REST Entities

| Entity External ID | Endpoint Pattern | Description |
|---|---|---|
| `SecretScanningAlert` | `/enterprises/{slug}/secret-scanning/alerts` or `/orgs/{orgName}/secret-scanning/alerts` | Secret scanning alerts |

## Notes

- Child entities with `$.*` external IDs are nested within their parent's GraphQL response.
- Entities with `uniqueId` as their unique attribute synthesize a composite ID from multiple fields.
- The adapter supports both Enterprise Cloud and Enterprise Server deployments.

## Adding New Entity External IDs

**Requires code changes.** Each entity has a dedicated GraphQL query builder and parsing logic. The adapter validates against `ValidEntityExternalIDs`.

### Adding a GraphQL entity:

1. Add a constant in `datasource.go` (e.g., `Project = "Project"`).
2. Add an entry to `ValidEntityExternalIDs` with the unique ID attribute, parse path, and any required attributes.
3. Create a new `QueryBuilder` struct in `query.go` with the GraphQL query for the new entity.
4. Add a case in `ConstructQueryBuilder` to dispatch to the new builder.
5. If it's a child entity (nested in a parent's response), set the appropriate `ParsePath` and parent reference.
6. Add tests.

### Adding a REST entity:

1. Add a constant in `datasource.go`.
2. Add the entity to `ValidEntityExternalIDs`.
3. Add endpoint mappings in `request.go` for both Enterprise Cloud and Enterprise Server deployment types.
4. Implement response parsing if the REST response format differs from existing entities.
5. Add tests.
