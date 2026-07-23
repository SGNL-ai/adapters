# LDAP Adapter

## Entity External IDs

The Entity External ID maps to an LDAP search filter query. Unlike REST-based adapters, external IDs are not used in URL construction - they serve as configuration keys that resolve to LDAP filter strings.

The connection target is a single LDAP server URL (e.g., `ldap://host:389` or `ldaps://host:636`).

| Entity External ID | Default LDAP Filter | Description |
|---|---|---|
| `User` | `(&(objectCategory=user)(objectClass=user)(distinguishedName=*))` | Active Directory users |
| `Group` | `(&(objectCategory=group)(objectClass=group)(distinguishedName=*))` | Active Directory groups |
| `Computer` | `(&(objectCategory=computer)(name=*))` | Active Directory computers |
| `GroupMember` | `(&(objectClass=group)({{CollectionAttribute}}={{CollectionId}}))` | Group members (child of Group) |

## How It Works

The Entity External ID is used as a lookup key into the `entityConfig` map to retrieve the LDAP filter query. That filter is then passed to an LDAP search operation against the configured `BaseDN`.

The adapter is open-ended - you can configure any entity external ID with a custom LDAP filter in the `entityConfig`:

```json
{
  "entityConfig": {
    "User": {"query": "(&(objectCategory=user)(objectClass=user)(distinguishedName=*))"},
    "Group": {"query": "(&(objectCategory=group)(objectClass=group)(distinguishedName=*))"},
    "Person": {"query": "(&(objectClass=person)(cn=*))"}
  }
}
```

## Parent-Child Relationships

- `GroupMember` is a child of `Group` - the adapter iterates groups and fetches their `member` attribute.
- Template placeholders `{{CollectionAttribute}}` and `{{CollectionId}}` in the query are replaced at runtime with the parent group's identifying attribute and value.

## Notes

- Default unique ID attributes use `distinguishedName`.
- Supports Active Directory-specific binary attributes (`objectGUID`, `objectSid`, etc.) with automatic hex conversion.
- Pagination uses LDAP Simple Paged Results Control.
- Max page size is 999.
- Two versions exist (v1.0.0 and v2.0.0) with different GroupMember query strategies.

## Adding New Entity External IDs

**No code changes required.** The adapter is driven entirely by the `entityConfig` map in the adapter configuration. To add a new entity, include it in the config with its LDAP filter query:

```json
{
  "entityConfig": {
    "OrganizationalUnit": {
      "query": "(&(objectCategory=organizationalUnit)(ou=*))"
    },
    "Contact": {
      "query": "(&(objectClass=contact)(mail=*))"
    }
  }
}
```

For a new **child** entity (scoped to a parent collection), also include `memberOf`, `collectionAttribute`, `memberAttribute`, and related fields:

```json
{
  "entityConfig": {
    "RoleMember": {
      "query": "(&(objectClass=group)({{CollectionAttribute}}={{CollectionId}}))",
      "memberOf": "Role",
      "collectionAttribute": "distinguishedName",
      "memberAttribute": "member",
      "memberUniqueIdAttribute": "distinguishedName",
      "memberOfUniqueIdAttribute": "distinguishedName"
    }
  }
}
```
