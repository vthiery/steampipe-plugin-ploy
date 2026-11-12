# Table: ploy_access_review

Access reviews: one per resource per campaign cycle, with reviewers and due date.

## Examples

**Open reviews past their due date:**

```sql
select name, status, due_date, resource_name
from ploy_access_review
where status in ('in_progress', 'ready_for_approval')
  and due_date::date < current_date;
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Review ID. |
| `name` | text | Name. |
| `status` | text | in_progress, ready_for_approval, complete or archived. |
| `campaign_id` | text | Campaign ID. |
| `cycle_id` | text | Cycle ID. |
| `resource_id` | text | Resource under review. |
| `resource_name` | text | Name of the resource under review. |
| `resource` | jsonb | Resource under review. |
| `reviewers` | jsonb | Assigned reviewers. |
| `due_date` | text | Effective deadline (YYYY-MM-DD). |
| `started_at` | timestamp with time zone | Start time. |
| `opened_at` | timestamp with time zone | When a reviewer first opened the review. |
| `completed_at` | timestamp with time zone | Completion time. |
| `excluded_at` | timestamp with time zone | When an admin removed the review from scope. |
| `created_at` | timestamp with time zone | Creation time. |
| `updated_at` | timestamp with time zone | Last update time. |
