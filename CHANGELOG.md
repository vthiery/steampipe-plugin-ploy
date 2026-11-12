## v0.0.2 [2026-11-12]

_What's new?_

- Added the `ploy_ip_whitelist` and `ploy_internal_domain` tables (organisation settings; need the `organisation:read` scope).

## v0.0.1 [2026-11-12]

_What's new?_

- Initial release with the following tables: `ploy_member`, `ploy_user`, `ploy_resource`, `ploy_resource_access`, `ploy_access_review_campaign`, `ploy_access_review`, `ploy_access_request`, `ploy_audit_log` and `ploy_integration`.
- `ploy_resource_access` lists access organisation-wide and pushes down the API's `memberStatus`, `memberIsManaged`, `active`, `shadow`, `identityActive`, `hasMember`, `identityType`, `resourceType`, `managedApp`, `accessorType`, `resourceId`, `memberId` and `identityId` filters. New columns: `shadow`, `identity_active` and the filter columns.
- Honour `Retry-After` and `X-RateLimit-Reset` on HTTP 429.
- Fixed filter-only columns of `ploy_resource_access` (`member_status`, `member_is_managed`, `has_member`, `identity_type`, `resource_type`, `managed_app`) returning no rows when an earlier query had cached the same data without the filter. They now use an exact query-cache match.
- Log the request URL at debug level.
