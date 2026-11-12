# Table: ploy_access_review_campaign

Access review campaigns (recurring user access review programmes).

## Examples

**List active campaigns:**

```sql
select name, review_type, next_cycle_date, auto_generate_cycles
from ploy_access_review_campaign
where status = 'active';
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Campaign ID. |
| `name` | text | Name. |
| `description` | text | Description. |
| `status` | text | active or archived. |
| `review_type` | text | standard, attestation_only or entitlement_only. |
| `auto_generate_cycles` | boolean | Whether new cycles are generated on schedule. |
| `due_days` | bigint | Days from cycle start until reviews are due. |
| `last_cycle_at` | timestamp with time zone | Last cycle start. |
| `next_cycle_date` | text | Next scheduled cycle (YYYY-MM-DD). |
| `created_at` | timestamp with time zone | Creation time. |
| `updated_at` | timestamp with time zone | Last update time. |
