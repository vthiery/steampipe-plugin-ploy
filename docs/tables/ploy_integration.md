# Table: ploy_integration

Integrations connected to Ploy (IdP, HRIS, SaaS, ...), with last scan time.

## Examples

**Integrations that have not scanned in 48 hours:**

```sql
select name, type, enabled, last_scanned
from ploy_integration
where last_scanned < now() - interval '48 hours';
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Integration ID. |
| `name` | text | Vendor slug, e.g. okta. |
| `connection_name` | text | Connection label. |
| `type` | text | idp, email, hris, scm, infra, saas, security, knowledge_source or custom. |
| `enabled` | boolean | Whether enabled. |
| `is_setup` | boolean | Whether setup is complete. |
| `last_scanned` | timestamp with time zone | Last successful scan. |
| `created_at` | timestamp with time zone | Creation time. |
| `updated_at` | timestamp with time zone | Last update time. |
