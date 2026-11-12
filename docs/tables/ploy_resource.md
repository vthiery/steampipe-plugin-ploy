# Table: ploy_resource

Resources (apps, groups, roles, devices, drives, AI agents, ...) known to Ploy, with owners and access config.

## Examples

**Count resources by type:**

```sql
select type, count(*) as total
from ploy_resource
group by type
order by total desc;
```

**Sanctioned apps without an owner:**

```sql
select name, domain
from ploy_resource
where type = 'app'
  and app_status = 'sanctioned'
  and jsonb_array_length(coalesce(owners, '[]'::jsonb)) = 0;
```

**Apps waiting for review:**

```sql
select name, domain, created_at
from ploy_resource
where type = 'app'
  and app_status = 'to_review'
order by created_at;
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Resource ID. |
| `name` | text | Name. |
| `type` | text | Resource type: app, group, installed_app, idp_app, device, drive, site, role, ai_agent, license, policy, account. |
| `status` | text | active or archived. |
| `app_status` | text | App review status: sanctioned, unsanctioned, to_review, in_review, not_recommended, ... |
| `is_managed_app` | boolean | Whether the app is managed in Ploy. |
| `domain` | text | App primary domain. |
| `ext_id` | text | External ID in the source system. |
| `integration_id` | text | Integration the resource was discovered from. |
| `description` | text | Description. |
| `owners` | jsonb | Owners (id, email, type primary/secondary). |
| `tags` | jsonb | Tags. |
| `attrs` | jsonb | Attributes. |
| `policy` | jsonb | Access config: approval policy, max duration, access conditions, low-usage checks. |
| `is_dynamic` | boolean | Membership controlled by dynamic rules in the source. |
| `dynamic_rules` | jsonb | Dynamic membership rules. |
| `deeplink` | text | Deep link. |
| `default_reviewer_instruction` | text | Default instruction for access reviewers. |
| `created_at` | timestamp with time zone | Creation time. |
| `updated_at` | timestamp with time zone | Last update time. |
