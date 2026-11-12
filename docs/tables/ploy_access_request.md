# Table: ploy_access_request

Access requests and their approval/provisioning lifecycle.

## Examples

**Requests still pending approval:**

```sql
select id, status, current_approval_stage, total_approval_stages, requested_at
from ploy_access_request
where status = 'pending'
order by requested_at;
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Request ID. |
| `status` | text | Lifecycle status. |
| `member_id` | text | Member receiving the access. |
| `member` | jsonb | Member receiving the access. |
| `requested_by_member` | jsonb | Member who initiated the request. |
| `requested_by_identity_id` | text | Identity that initiated an agent/NHI request. |
| `selectable_resource` | jsonb | Selectable resource requested. |
| `lines` | jsonb | Requested roles and side effects. |
| `pending_approvers` | jsonb | Approvers the current stage waits on. |
| `current_approval_stage` | bigint | Current approval stage. |
| `total_approval_stages` | bigint | Number of approval stages. |
| `requested_indefinite` | boolean | Requested with no end date. |
| `requested_for_seconds` | bigint | Requested duration. |
| `source` | text | member_portal, browser_extension, slack or api. |
| `comments` | text | Requester comments. |
| `requested_at` | timestamp with time zone | Request time. |
| `approval_decision_made_at` | timestamp with time zone | Decision time. |
| `provisioned_at` | timestamp with time zone | Provisioning time. |
| `created_at` | timestamp with time zone | Creation time. |
| `updated_at` | timestamp with time zone | Last update time. |
