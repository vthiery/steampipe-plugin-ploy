# Ploy Plugin for Steampipe

Use SQL to query members, resources, access, access reviews, access requests and audit logs from [Ploy](https://joinploy.com), the identity governance and access management platform.

- **[Get started →](https://github.com/vthiery/steampipe-plugin-ploy)**
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-ploy/issues)

## Installation

```sh
steampipe plugin install ghcr.io/vthiery/ploy
```

## Configuration

Copy the sample config:

```sh
cp config/ploy.spc ~/.steampipe/config/ploy.spc
```

Edit `~/.steampipe/config/ploy.spc`:

```hcl
connection "ploy" {
  plugin = "ghcr.io/vthiery/ploy"

  # API key for the Ploy public API. Read-only scopes are sufficient.
  # Falls back to the PLOY_API_KEY environment variable when not set.
  api_key = "YOUR_API_KEY"

  # Base URL of the API. Defaults to https://api.joinploy.com/api/v1.
  # base_url = "https://api.joinploy.com/api/v1"
}
```

The scopes the key needs depend on the tables you query: `members:read` for `ploy_member`, `users:read` for `ploy_user`, `resources:read` for `ploy_resource`, `resource_access:read` for `ploy_resource_access`, `access_reviews:read` for the two access review tables, `access_requests:read` for `ploy_access_request`, `audit_logs:read` for `ploy_audit_log` and `integrations:read` for `ploy_integration`.

## Tables

| Table | Description |
|-------|-------------|
| [ploy_member](tables/ploy_member.md) | Org members (people) with HR attributes and their identities. |
| [ploy_user](tables/ploy_user.md) | Users of the Ploy dashboard (owners and admins). |
| [ploy_resource](tables/ploy_resource.md) | Apps, groups, roles, devices, drives, AI agents and more, with owners and access config. |
| [ploy_resource_access](tables/ploy_resource_access.md) | Who holds access to which resource. |
| [ploy_access_review_campaign](tables/ploy_access_review_campaign.md) | Access review campaigns. |
| [ploy_access_review](tables/ploy_access_review.md) | Access reviews, one per resource per campaign cycle. |
| [ploy_access_request](tables/ploy_access_request.md) | Access requests and their approval and provisioning lifecycle. |
| [ploy_audit_log](tables/ploy_audit_log.md) | The organisation audit trail. |
| [ploy_integration](tables/ploy_integration.md) | Connected integrations (identity provider, HRIS, SaaS, ...) with last scan time. |

## Rate limiting

Ploy does not publish a rate limit. The plugin sends at most 8 requests per second and retries HTTP 429 and 5xx responses with exponential backoff (0.5 s up to 30 s, 8 retries).
