# Table: ploy_resource_access

Access held to resources, organisation-wide.

Every filter is optional and is pushed down to the API. Use the filter columns to keep scans small: the API refuses a query that matches more than 250,000 rows unless it filters on `resource_id`, `member_id` or `identity_id`. Filter columns (`member_status`, `member_is_managed`, `has_member`, `identity_type`, `resource_type`, `managed_app`) return the value you filtered on, and null otherwise.

## Examples

**Leavers whose accounts are still enabled:**

```sql
select
  member_id,
  resource_id,
  granted_at
from
  ploy_resource_access
where
  member_status = 'inactive'
  and member_is_managed
  and active
  and identity_active;
```

**Access held by accounts not matched to a person:**

```sql
select
  identity_id,
  resource_id,
  identity_active
from
  ploy_resource_access
where
  has_member = false
  and active;
```

**Shadow access (seen in use, not confirmed by a source of truth):**

```sql
select
  resource_id,
  member_id,
  granted_at
from
  ploy_resource_access
where
  shadow
  and active;
```

**Active access held by one member:**

```sql
select
  resource_id,
  granted_at,
  expires_at
from
  ploy_resource_access
where
  member_id = '<member id>'
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
| `shadow` | boolean | Access Ploy saw in use that no source of truth has confirmed (Shadow Access). |
| `identity_active` | boolean | Whether the holding account is enabled in its own system. False once deactivated or removed there, e.g. a user deprovisioned in Okta. |
| `managed_status` | text | Ploy provisioning lifecycle state. |
| `granted_at` | timestamp with time zone | Grant time. |
| `expires_at` | timestamp with time zone | Expiry for time-boxed access. |
| `revoked_at` | timestamp with time zone | Revocation time. |
| `suspended_at` | timestamp with time zone | Suspension time. |
| `grants` | jsonb | Entitlements held through this access. |
| `created_at` | timestamp with time zone | Creation time. |
| `member_status` | text | Status of the member holding the access: active, inactive or inactive_onboarding. Filter only: null unless used in the WHERE clause. |
| `member_is_managed` | boolean | Whether the holding member is managed. With member_status = 'inactive' this selects leavers (terminated members). Filter only: null unless used in the WHERE clause. |
| `has_member` | boolean | true: identities owned by a member. false: service accounts and accounts not matched to a person. Filter only: null unless used in the WHERE clause. |
| `identity_type` | text | Identity type: human, service_account, service_principal, managed_identity, app_registration, api_key, agent_identity or shared_mailbox. Filter only: null unless used in the WHERE clause. |
| `resource_type` | text | Type of the resource the access is to, as type on ploy_resource. Filter only: null unless used in the WHERE clause. |
| `managed_app` | boolean | true: access to apps on the Managed apps list. false: every other resource. Filter only: null unless used in the WHERE clause. |
