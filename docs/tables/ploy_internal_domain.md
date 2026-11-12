# Table: ploy_internal_domain

Email domains the organisation declared as internal. Members on these domains count as employees rather than outside collaborators; with none declared, every member counts as internal.

## Examples

**List the internal domains:**

```sql
select
  domain
from
  ploy_internal_domain
order by
  domain;
```

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `domain` | text | Internal email domain, without the @. |
