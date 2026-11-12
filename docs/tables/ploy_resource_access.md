# Table: ploy_resource_access

Access held to resources. Filter on resource_id, member_id or identity_id for a direct lookup; without one, the table fans out over every member (one API call each), so resource-held access is not listed.

Filter on `resource_id`, `member_id` or `identity_id` for a direct lookup. Without one of them the table lists access for every member, which costs one API call per member and is slow on large organisations.

## Examples

**Active access held by one member:**

```sql
select resource_id, granted_at, expires_at, managed_status
from ploy_resource_access
where member_id = '<member id>'
  and active;
```

**Who has access to one resource:**

```sql
select member_id, identity_id, granted_at
from ploy_resource_access
where resource_id = '<resource id>'
  and active;
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Resource access ID. |
| `resource_id` | text | Resource the access is to. |
| `member_id` | text | Member owning the holding identity. |
| `identity_id` | text | Identity holding the access. |
| `accessor_type` | text | identity or resource. |
| `accessing_resource_id` | text | Resource holding the access (accessor_type = resource). |
| `active` | boolean | Neither revoked nor suspended. |
| `managed_status` | text | Ploy provisioning lifecycle state. |
| `granted_at` | timestamp with time zone | Grant time. |
| `expires_at` | timestamp with time zone | Expiry for time-boxed access. |
| `revoked_at` | timestamp with time zone | Revocation time. |
| `suspended_at` | timestamp with time zone | Suspension time. |
| `grants` | jsonb | Entitlements held through this access. |
| `created_at` | timestamp with time zone | Creation time. |
