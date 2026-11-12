# Table: ploy_audit_log

The organisation's Ploy audit trail. Filter occurred_at (>, >=, <, <=) to bound the scan.

Filter on `occurred_at` to bound the scan: the comparison is sent to the API as `occurredSince` / `occurredBefore`.

## Examples

**Recent audit events:**

```sql
select occurred_at, type, actor_type, actor_email, action_result
from ploy_audit_log
where occurred_at > now() - interval '7 days'
order by occurred_at desc
limit 50;
```

**Failed actions in the last 30 days:**

```sql
select occurred_at, type, actor_email
from ploy_audit_log
where occurred_at > now() - interval '30 days'
  and action_result = 'failed';
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Entry ID. |
| `type` | text | Event type, e.g. workflows.executed.create_resource. |
| `occurred_at` | timestamp with time zone | When the action started. |
| `action_result` | text | success or failed; null for events with no outcome. |
| `actor_type` | text | user, member, api_key or system. |
| `actor_user_id` | text | Acting user ID. |
| `actor_member_id` | text | Acting member ID. |
| `actor_api_key_id` | text | Acting API key ID. |
| `actor_email` | text | Acting user's/member's email. |
| `actor_name` | text | Acting user's/member's name (`Ploy Staff` for Ploy support). |
| `actor` | jsonb | Actor object. |
| `ip_address` | text | Source IP. |
| `user_agent` | text | Source user agent. |
| `payload` | jsonb | Event details (secrets redacted). |
