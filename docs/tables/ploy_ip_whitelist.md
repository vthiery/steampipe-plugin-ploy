# Table: ploy_ip_whitelist

IP allow-lists of the organisation, one row per surface: the admin dashboard and the employee portal (which also gates the browser extension). An empty list means the surface is unrestricted, which is the default.

## Examples

**Surfaces without an IP restriction:**

```sql
select
  type
from
  ploy_ip_whitelist
where
  jsonb_array_length(cidrs) = 0;
```

**Allowed ranges of the admin dashboard:**

```sql
select
  jsonb_array_elements_text(cidrs) as cidr
from
  ploy_ip_whitelist
where
  type = 'dashboard';
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `type` | text | Surface the list gates: dashboard or employee-portal. |
| `cidrs` | jsonb | Allowed ranges in CIDR notation. Empty means unrestricted. |
