# Table: ploy_user

Users of the Ploy dashboard (who administers Ploy itself).

## Examples

**List Ploy dashboard administrators:**

```sql
select email, role, status
from ploy_user
where role in ('owner', 'admin');
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | User ID. |
| `email` | text | Email. |
| `first_name` | text | First name. |
| `last_name` | text | Last name. |
| `role` | text | Dashboard role: owner, admin, member, ... |
| `status` | text | Org-user status (active, invited, ...). |
| `permissions` | jsonb | Granted permissions. |
| `resource_allow_list` | jsonb | Resources this user is restricted to. |
| `resource_block_list` | jsonb | Resources this user is blocked from. |
| `created_at` | timestamp with time zone | Creation time. |
