package ploy

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Every filter is optional and pushed down to GET /resource-access. A query matching more
// than 250,000 rows without a resource_id, member_id or identity_id is refused by the API.
var resourceAccessQuals = []qualParam{
	{"resource_id", "resourceId"},
	{"member_id", "memberId"},
	{"identity_id", "identityId"},
	{"active", "active"},
	{"shadow", "shadow"},
	{"accessor_type", "accessorType"},
	{"member_status", "memberStatus"},
	{"member_is_managed", "memberIsManaged"},
	{"has_member", "hasMember"},
	{"identity_active", "identityActive"},
	{"identity_type", "identityType"},
	{"resource_type", "resourceType"},
	{"managed_app", "managedApp"},
}

// filterCol builds a column that only exists to push a filter to the API. Its value is the
// value the query filtered on (copied onto each row by the list hydrate), and null when the
// query did not filter on it.
func filterCol(name string, t proto.ColumnType, field, desc string) *plugin.Column {
	return col(name, t, field, desc+" Filter only: null unless used in the WHERE clause.")
}

func tablePloyResourceAccess() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_resource_access",
		Description: "Access held to resources, organisation-wide. Filter to keep scans small; see the filter columns.",
		List: &plugin.ListConfig{
			Hydrate:    listHydrate("/resource-access", resourceAccessQuals),
			KeyColumns: optionalQuals(resourceAccessQuals, "active", "shadow", "member_is_managed", "has_member", "identity_active", "managed_app"),
		},
		Get: &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/resource-access")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Resource access ID."),
			col("resource_id", tStr, "resourceId", "Resource the access is to."),
			col("member_id", tStr, "memberId", "Member owning the holding identity."),
			col("identity_id", tStr, "identityId", "Identity holding the access."),
			col("accessor_type", tStr, "accessorType", "identity or resource."),
			col("accessing_resource_id", tStr, "accessingResourceId", "Resource holding the access (accessor_type = resource)."),
			col("active", tBool, "active", "Neither revoked nor suspended."),
			col("shadow", tBool, "shadow", "Access Ploy saw in use that no source of truth has confirmed (Shadow Access)."),
			col("identity_active", tBool, "identityActive", "Whether the holding account is enabled in its own system. False once deactivated or removed there, e.g. a user deprovisioned in Okta."),
			col("managed_status", tStr, "managedStatus", "Ploy provisioning lifecycle state."),
			col("granted_at", tTS, "grantedAt", "Grant time."),
			col("expires_at", tTS, "expiresAt", "Expiry for time-boxed access."),
			col("revoked_at", tTS, "revokedAt", "Revocation time."),
			col("suspended_at", tTS, "suspendedAt", "Suspension time."),
			col("grants", tJSON, "grants", "Entitlements held through this access."),
			col("created_at", tTS, "createdAt", "Creation time."),
			filterCol("member_status", tStr, "memberStatus", "Status of the member holding the access: active, inactive or inactive_onboarding."),
			filterCol("member_is_managed", tBool, "memberIsManaged", "Whether the holding member is managed. With member_status = 'inactive' this selects leavers (terminated members)."),
			filterCol("has_member", tBool, "hasMember", "true: identities owned by a member. false: service accounts and accounts not matched to a person."),
			filterCol("identity_type", tStr, "identityType", "Identity type: human, service_account, service_principal, managed_identity, app_registration, api_key, agent_identity or shared_mailbox."),
			filterCol("resource_type", tStr, "resourceType", "Type of the resource the access is to, as type on ploy_resource."),
			filterCol("managed_app", tBool, "managedApp", "true: access to apps on the Managed apps list. false: every other resource."),
		},
	}
}
