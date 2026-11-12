# Table: ploy_member

Org members (people) known to Ploy, with HR attributes and their non-app identities.

## Examples

**List active members:**

```sql
select email, first_name, last_name, department, job_title
from ploy_member
where status = 'active'
order by email;
```

**Members that are inactive but still have an IdP status:**

```sql
select email, status, source_of_truth_status
from ploy_member
where status = 'inactive'
  and source_of_truth_status is not null;
```

**Active members without a manager:**

```sql
select email, department, job_title
from ploy_member
where status = 'active'
  and manager_id is null;
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `id` | text | Member ID. |
| `email` | text | Primary email. |
| `first_name` | text | First name. |
| `last_name` | text | Last name. |
| `status` | text | active, inactive or inactive_onboarding. |
| `internal` | boolean | True for internal (employee) members. |
| `source_of_truth_status` | text | Status reported by the member's source of truth (e.g. Okta ACTIVE/SUSPENDED/DEPROVISIONED). |
| `is_managed` | boolean | Whether the member is managed by Ploy. |
| `is_manager` | boolean | Whether the member manages others. |
| `manager_id` | text | Manager's member ID. |
| `manager` | jsonb | Manager reference. |
| `employee_id` | text | HR employee ID. |
| `employment_type` | text | Employment type. |
| `employment_start_date` | text | Employment start date (YYYY-MM-DD). |
| `employment_end_date` | text | Employment end date (YYYY-MM-DD). |
| `last_working_day` | text | Last working day (YYYY-MM-DD). |
| `job_title` | text | Job title. |
| `department` | text | Department. |
| `team` | text | Team. |
| `business_unit` | text | Business unit. |
| `role_level` | text | Role level. |
| `location_city` | text | City. |
| `location_country` | text | Country. |
| `location_region` | text | Region. |
| `working_location` | text | Working location. |
| `identities` | jsonb | Active non-app identities (type, authType, username, active). |
| `custom_fields` | jsonb | Custom field values. |
| `locked_fields` | jsonb | Fields locked against sync. |
| `created_at` | timestamp with time zone | Creation time. |
| `updated_at` | timestamp with time zone | Last update time. |
