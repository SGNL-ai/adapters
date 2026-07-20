# Azure AD (Entra ID) Adapter

## Entity External IDs

The Entity External ID maps to a Microsoft Graph API resource path. The adapter constructs URLs as:

```
{baseURL}/{apiVersion}/{path}
```

Where `apiVersion` is `v1.0`.

| Entity External ID | API Path | Description |
|---|---|---|
| `User` | `/users` | Azure AD users |
| `Group` | `/groups` | Azure AD groups |
| `GroupMember` | `/groups/{groupID}/members` | Members of a specific group (child of Group) |
| `Application` | `/applications` | Registered applications |
| `Device` | `/devices` | Registered devices |
| `Role` | `/directoryRoles` | Directory roles |
| `RoleMember` | `/users/{userID}/transitiveMemberOf/microsoft.graph.directoryRole` | Roles assigned to a user (child of User) |
| `RoleAssignment` | `/roleManagement/directory/roleAssignments` | Role assignment records |
| `RoleAssignmentScheduleRequest` | `/roleManagement/directory/roleAssignmentScheduleRequests` | PIM role assignment schedule requests |
| `GroupAssignmentScheduleRequest` | `/identityGovernance/privilegedAccess/group/assignmentScheduleRequests` | PIM group assignment schedule requests |

## Parent-Child Relationships

- `GroupMember` is a child of `Group` - the adapter iterates groups and fetches members per group.
- `RoleMember` is a child of `User` - the adapter iterates users and fetches their transitive role memberships.

## Adding New Entity External IDs

**Requires code changes.** The adapter validates against a fixed `ValidEntityExternalIDs` map and uses a switch statement in `endpoint.go` to construct entity-specific URL paths.

To add a new entity:

1. Add a constant in `datasource.go` (e.g., `ServicePrincipal string = "ServicePrincipal"`).
2. Add an entry to `ValidEntityExternalIDs` in `datasource.go` specifying whether it has a parent (`memberOf`).
3. Add a case in `ConstructEndpoint` in `endpoint.go` with the Microsoft Graph API path (e.g., `/servicePrincipals`).
4. If the entity is a child (member), configure `memberOf` to point to the parent entity and implement collection-based pagination.
5. If the entity needs `$expand` support, add it to the `expandableRelationships` map.
6. Add tests in `endpoint_test.go` and `adapter_test.go`.

## Notes

- GroupMember supports type-filtering via advanced filters (e.g., only user members via `/microsoft.graph.user` suffix).
- Expandable relationships are supported on User, Group, Application, and Device entities (e.g., `$expand=manager`).
- PIM entities use `$top`/`$skip` pagination instead of `@odata.nextLink` cursors.
