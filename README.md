# Ploy Plugin for Steampipe

Use SQL to query members, resources, access, access reviews, access requests and audit logs from [Ploy](https://joinploy.com), the identity governance and access management platform.

- **[Get started →](https://github.com/vthiery/steampipe-plugin-ploy)**
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-ploy/issues)

## Quick start

### Install

Install the plugin with [Steampipe](https://steampipe.io):

```sh
steampipe plugin install ghcr.io/vthiery/ploy
```

### Configure

Copy the sample config and set your API key:

```sh
cp config/ploy.spc ~/.steampipe/config/ploy.spc
```

Edit `~/.steampipe/config/ploy.spc`:

```hcl
connection "ploy" {
  plugin = "ghcr.io/vthiery/ploy"

  # API key for the Ploy public API. Read-only scopes are sufficient.
  api_key = "YOUR_API_KEY"
}
```

Alternatively, leave `api_key` unset and export `PLOY_API_KEY`.

### Run a query

```shell
steampipe query
```

List active members without a manager:

```sql
select
  email,
  department,
  job_title
from
  ploy_member
where
  status = 'active'
  and manager_id is null;
```

## Tables

| Table | Description |
|-------|-------------|
| [ploy_member](docs/tables/ploy_member.md) | Org members (people) with HR attributes and their identities. |
| [ploy_user](docs/tables/ploy_user.md) | Users of the Ploy dashboard (owners and admins). |
| [ploy_resource](docs/tables/ploy_resource.md) | Apps, groups, roles, devices, drives, AI agents and more, with owners and access config. |
| [ploy_resource_access](docs/tables/ploy_resource_access.md) | Who holds access to which resource. |
| [ploy_access_review_campaign](docs/tables/ploy_access_review_campaign.md) | Access review campaigns. |
| [ploy_access_review](docs/tables/ploy_access_review.md) | Access reviews, one per resource per campaign cycle. |
| [ploy_access_request](docs/tables/ploy_access_request.md) | Access requests and their approval and provisioning lifecycle. |
| [ploy_audit_log](docs/tables/ploy_audit_log.md) | The organisation audit trail. |
| [ploy_integration](docs/tables/ploy_integration.md) | Connected integrations (identity provider, HRIS, SaaS, ...) with last scan time. |

> **Note — API limits:** Ploy does not publish a rate limit. The plugin sends at most 8 requests per second and retries on HTTP 429 with exponential backoff. `ploy_resource_access` can only be listed per member, resource or identity, so an unfiltered scan makes one API call per member; filter on `resource_id`, `member_id` or `identity_id` whenever you can. When joining it to other tables, read the other table in a `MATERIALIZED` CTE and reach `ploy_resource_access` with `CROSS JOIN LATERAL`, otherwise the query planner can turn the join into thousands of API calls.

## Development

### Prerequisites

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

### Build and Install

```sh
make install
```

Configure the plugin as a local build:

```sh
cp config/ploy.spc ~/.steampipe/config/ploy.spc
```

and set `plugin = "local/ploy"` in the connection.

## Testing

Run a smoke query against every table:

```sh
make test
```

The test script ([scripts/test_tables.sh](scripts/test_tables.sh)) builds the plugin, queries each table, and reports pass/fail/skip.

### Further reading

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Ploy API reference](https://app.joinploy.com/api-docs)
