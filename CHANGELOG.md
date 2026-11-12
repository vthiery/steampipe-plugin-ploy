## v0.0.1 [2026-11-12]

_What's new?_

- Initial release with the following tables: `ploy_member`, `ploy_user`, `ploy_resource`, `ploy_resource_access`, `ploy_access_review_campaign`, `ploy_access_review`, `ploy_access_request`, `ploy_audit_log` and `ploy_integration`.
- `ploy_resource_access` lists access organisation-wide and pushes down the API's `memberStatus`, `memberIsManaged`, `active`, `shadow`, `identityActive`, `hasMember`, `identityType`, `resourceType`, `managedApp`, `accessorType`, `resourceId`, `memberId` and `identityId` filters. New columns: `shadow`, `identity_active` and the filter columns.
- Honour `Retry-After` and `X-RateLimit-Reset` on HTTP 429.
