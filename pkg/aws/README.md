# AWS IAM Adapter

## Entity External IDs

The Entity External ID maps to AWS IAM SDK API calls. Unlike REST-based adapters, this adapter uses the AWS SDK v2 rather than constructing HTTP URLs directly.

| Entity External ID | AWS IAM API Call (List) | AWS IAM API Call (Get) | Description |
|---|---|---|---|
| `User` | `ListUsers` | `GetUser` | IAM users |
| `Group` | `ListGroups` | `GetGroup` | IAM groups |
| `GroupMember` | `GetGroup` (returns `.Users`) | - | Group members (child of Group) |
| `Role` | `ListRoles` | `GetRole` | IAM roles |
| `IdentityProvider` | `ListSAMLProviders` | - | SAML identity providers |
| `Policy` | `ListPolicies` | `GetPolicy` | IAM policies |
| `RolePolicy` | `ListAttachedRolePolicies` | - | Policies attached to roles (child of Role) |
| `UserPolicy` | `ListAttachedUserPolicies` | - | Policies attached to users (child of User) |
| `GroupPolicy` | `ListAttachedGroupPolicies` | - | Policies attached to groups (child of Group) |

## Parent-Child Relationships

| Child Entity | Parent Entity | Collection Attribute |
|---|---|---|
| `GroupMember` | `Group` | `GroupName` |
| `GroupPolicy` | `Group` | `GroupName` |
| `RolePolicy` | `Role` | `RoleName` |
| `UserPolicy` | `User` | `UserName` |

Child entities iterate through their parent collection and fetch associated records per parent.

## Unique ID Construction

- Top-level entities use an `id` field constructed from their ARN.
- Member/policy entities use a composite ID: `{memberUniqueID}-{parentUniqueID}`.

## Notes

- Supports cross-account access via STS `AssumeRole` with configured role ARNs.
- Max page size is 1000, max resource accounts is 100.
- The adapter uses the AWS SDK internally - the underlying endpoint is `https://iam.amazonaws.com/`.

## Adding New Entity External IDs

**Requires code changes.** The adapter validates against `ValidEntityExternalIDs` and uses a switch statement in `datasource.go` to dispatch to entity-specific SDK call handlers.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `InstanceProfile string = "InstanceProfile"`).
2. Add an entry to `ValidEntityExternalIDs` with the ARN attribute, unique name field, and parent relationship (if any).
3. Create a handler struct in `iam.go` implementing the list/get SDK calls (following the pattern of `UserHandler`, `GroupHandler`, etc.).
4. Add a case in the switch statement in `GetPage` to dispatch to the new handler via `FetchEntities`.
5. If it's a child entity, set `MemberOf` and `CollectionAttribute` to configure parent iteration.
6. Add tests.
