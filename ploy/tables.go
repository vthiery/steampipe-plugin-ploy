package ploy

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const (
	tStr  = proto.ColumnType_STRING
	tBool = proto.ColumnType_BOOL
	tInt  = proto.ColumnType_INT
	tTS   = proto.ColumnType_TIMESTAMP
	tJSON = proto.ColumnType_JSON
)

//// ploy_member — people in the org (HRIS/IdP-sourced), with their identities.

var memberQuals = []qualParam{{"email", "email"}, {"status", "status"}}

func tablePloyMember() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_member",
		Description: "Org members (people) known to Ploy, with HR attributes and their non-app identities.",
		List:        &plugin.ListConfig{Hydrate: listHydrate("/members", memberQuals), KeyColumns: optionalQuals(memberQuals)},
		Get:         &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/members")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Member ID."),
			col("email", tStr, "email", "Primary email."),
			col("first_name", tStr, "firstName", "First name."),
			col("last_name", tStr, "lastName", "Last name."),
			col("status", tStr, "status", "active, inactive or inactive_onboarding."),
			col("internal", tBool, "internal", "True for internal (employee) members."),
			col("source_of_truth_status", tStr, "sourceOfTruthStatus", "Status reported by the member's source of truth (e.g. Okta ACTIVE/SUSPENDED/DEPROVISIONED)."),
			col("is_managed", tBool, "isManaged", "Whether the member is managed by Ploy."),
			col("is_manager", tBool, "isManager", "Whether the member manages others."),
			col("manager_id", tStr, "managerId", "Manager's member ID."),
			col("manager", tJSON, "manager", "Manager reference."),
			col("employee_id", tStr, "employeeId", "HR employee ID."),
			col("employment_type", tStr, "employmentType", "Employment type."),
			col("employment_start_date", tStr, "employmentStartDate", "Employment start date (YYYY-MM-DD)."),
			col("employment_end_date", tStr, "employmentEndDate", "Employment end date (YYYY-MM-DD)."),
			col("last_working_day", tStr, "lastWorkingDay", "Last working day (YYYY-MM-DD)."),
			col("job_title", tStr, "jobTitle", "Job title."),
			col("department", tStr, "department", "Department."),
			col("team", tStr, "team", "Team."),
			col("business_unit", tStr, "businessUnit", "Business unit."),
			col("role_level", tStr, "roleLevel", "Role level."),
			col("location_city", tStr, "locationCity", "City."),
			col("location_country", tStr, "locationCountry", "Country."),
			col("location_region", tStr, "locationRegion", "Region."),
			col("working_location", tStr, "workingLocation", "Working location."),
			col("identities", tJSON, "identities", "Active non-app identities (type, authType, username, active)."),
			col("custom_fields", tJSON, "customFields", "Custom field values."),
			col("locked_fields", tJSON, "lockedFields", "Fields locked against sync."),
			col("created_at", tTS, "createdAt", "Creation time."),
			col("updated_at", tTS, "updatedAt", "Last update time."),
		},
	}
}

//// ploy_user — Ploy dashboard users (admins/owners of the Ploy tenant itself).

var userQuals = []qualParam{{"email", "email"}}

func tablePloyUser() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_user",
		Description: "Users of the Ploy dashboard (who administers Ploy itself).",
		List:        &plugin.ListConfig{Hydrate: listHydrate("/users", userQuals), KeyColumns: optionalQuals(userQuals)},
		Get:         &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/users")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "User ID."),
			col("email", tStr, "email", "Email."),
			col("first_name", tStr, "firstName", "First name."),
			col("last_name", tStr, "lastName", "Last name."),
			col("role", tStr, "role", "Dashboard role: owner, admin, member, ..."),
			col("status", tStr, "status", "Org-user status (active, invited, ...)."),
			col("permissions", tJSON, "permissions", "Granted permissions."),
			col("resource_allow_list", tJSON, "resourceAllowList", "Resources this user is restricted to."),
			col("resource_block_list", tJSON, "resourceBlockList", "Resources this user is blocked from."),
			col("created_at", tTS, "createdAt", "Creation time."),
		},
	}
}

//// ploy_resource — apps, groups, roles, devices, ... discovered or managed by Ploy.

var resourceQuals = []qualParam{
	{"type", "type"}, {"status", "status"}, {"name", "name"}, {"ext_id", "extId"},
	{"domain", "domain"}, {"app_status", "appStatus"}, {"integration_id", "integrationId"},
}

func tablePloyResource() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_resource",
		Description: "Resources (apps, groups, roles, devices, drives, AI agents, ...) known to Ploy, with owners and access config.",
		List:        &plugin.ListConfig{Hydrate: listHydrate("/resources", resourceQuals), KeyColumns: optionalQuals(resourceQuals)},
		Get:         &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/resources")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Resource ID."),
			col("name", tStr, "name", "Name."),
			col("type", tStr, "type", "Resource type: app, group, installed_app, idp_app, device, drive, site, role, ai_agent, license, policy, account."),
			col("status", tStr, "status", "active or archived."),
			col("app_status", tStr, "appStatus", "App review status: sanctioned, unsanctioned, to_review, in_review, not_recommended, ..."),
			col("is_managed_app", tBool, "isManagedApp", "Whether the app is managed in Ploy."),
			col("domain", tStr, "domain", "App primary domain."),
			col("ext_id", tStr, "extId", "External ID in the source system."),
			col("integration_id", tStr, "integrationId", "Integration the resource was discovered from."),
			col("description", tStr, "description", "Description."),
			col("owners", tJSON, "owners", "Owners (id, email, type primary/secondary)."),
			col("tags", tJSON, "tags", "Tags."),
			col("attrs", tJSON, "attrs", "Attributes."),
			col("policy", tJSON, "policy", "Access config: approval policy, max duration, access conditions, low-usage checks."),
			col("is_dynamic", tBool, "isDynamic", "Membership controlled by dynamic rules in the source."),
			col("dynamic_rules", tJSON, "dynamicRules", "Dynamic membership rules."),
			col("deeplink", tStr, "deeplink", "Deep link."),
			col("default_reviewer_instruction", tStr, "defaultReviewerInstruction", "Default instruction for access reviewers."),
			col("created_at", tTS, "createdAt", "Creation time."),
			col("updated_at", tTS, "updatedAt", "Last update time."),
		},
	}
}

//// ploy_resource_access — who holds access to what.
// See table_ploy_resource_access.go for the list hydrate.

func tablePloyResourceAccess() *plugin.Table {
	return &plugin.Table{
		Name: "ploy_resource_access",
		Description: "Access held to resources. Filter on resource_id, member_id or identity_id for a direct lookup; " +
			"without one, the table fans out over every member (one API call each), so resource-held access is not listed.",
		List: &plugin.ListConfig{
			Hydrate:    listResourceAccess,
			KeyColumns: optionalQuals(resourceAccessQuals),
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
			col("managed_status", tStr, "managedStatus", "Ploy provisioning lifecycle state."),
			col("granted_at", tTS, "grantedAt", "Grant time."),
			col("expires_at", tTS, "expiresAt", "Expiry for time-boxed access."),
			col("revoked_at", tTS, "revokedAt", "Revocation time."),
			col("suspended_at", tTS, "suspendedAt", "Suspension time."),
			col("grants", tJSON, "grants", "Entitlements held through this access."),
			col("created_at", tTS, "createdAt", "Creation time."),
		},
	}
}

//// ploy_access_review_campaign

var campaignQuals = []qualParam{{"status", "status"}}

func tablePloyAccessReviewCampaign() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_access_review_campaign",
		Description: "Access review campaigns (recurring user access review programmes).",
		List:        &plugin.ListConfig{Hydrate: listHydrate("/access-review-campaigns", campaignQuals), KeyColumns: optionalQuals(campaignQuals)},
		Get:         &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/access-review-campaigns")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Campaign ID."),
			col("name", tStr, "name", "Name."),
			col("description", tStr, "description", "Description."),
			col("status", tStr, "status", "active or archived."),
			col("review_type", tStr, "reviewType", "standard, attestation_only or entitlement_only."),
			col("auto_generate_cycles", tBool, "autoGenerateCycles", "Whether new cycles are generated on schedule."),
			col("due_days", tInt, "dueDays", "Days from cycle start until reviews are due."),
			col("last_cycle_at", tTS, "lastCycleAt", "Last cycle start."),
			col("next_cycle_date", tStr, "nextCycleDate", "Next scheduled cycle (YYYY-MM-DD)."),
			col("created_at", tTS, "createdAt", "Creation time."),
			col("updated_at", tTS, "updatedAt", "Last update time."),
		},
	}
}

//// ploy_access_review — one review per resource per cycle.

var reviewQuals = []qualParam{
	{"status", "status"}, {"campaign_id", "campaignId"}, {"cycle_id", "cycleId"}, {"resource_id", "resourceId"},
}

func tablePloyAccessReview() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_access_review",
		Description: "Access reviews: one per resource per campaign cycle, with reviewers and due date.",
		List:        &plugin.ListConfig{Hydrate: listHydrate("/access-reviews", reviewQuals), KeyColumns: optionalQuals(reviewQuals)},
		Get:         &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/access-reviews")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Review ID."),
			col("name", tStr, "name", "Name."),
			col("status", tStr, "status", "in_progress, ready_for_approval, complete or archived."),
			col("campaign_id", tStr, "campaignId", "Campaign ID."),
			col("cycle_id", tStr, "cycleId", "Cycle ID."),
			col("resource_id", tStr, "resource.id", "Resource under review."),
			col("resource_name", tStr, "resource.name", "Name of the resource under review."),
			col("resource", tJSON, "resource", "Resource under review."),
			col("reviewers", tJSON, "reviewers", "Assigned reviewers."),
			col("due_date", tStr, "dueDate", "Effective deadline (YYYY-MM-DD)."),
			col("started_at", tTS, "startedAt", "Start time."),
			col("opened_at", tTS, "openedAt", "When a reviewer first opened the review."),
			col("completed_at", tTS, "completedAt", "Completion time."),
			col("excluded_at", tTS, "excludedAt", "When an admin removed the review from scope."),
			col("created_at", tTS, "createdAt", "Creation time."),
			col("updated_at", tTS, "updatedAt", "Last update time."),
		},
	}
}

//// ploy_access_request

var requestQuals = []qualParam{{"member_id", "memberId"}, {"status", "status"}}

func tablePloyAccessRequest() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_access_request",
		Description: "Access requests and their approval/provisioning lifecycle.",
		List:        &plugin.ListConfig{Hydrate: listHydrate("/access-requests", requestQuals), KeyColumns: optionalQuals(requestQuals)},
		Get:         &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/access-requests")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Request ID."),
			col("status", tStr, "status", "Lifecycle status."),
			col("member_id", tStr, "member.id", "Member receiving the access."),
			col("member", tJSON, "member", "Member receiving the access."),
			col("requested_by_member", tJSON, "requestedByMember", "Member who initiated the request."),
			col("requested_by_identity_id", tStr, "requestedByIdentityId", "Identity that initiated an agent/NHI request."),
			col("selectable_resource", tJSON, "selectableResource", "Selectable resource requested."),
			col("lines", tJSON, "lines", "Requested roles and side effects."),
			col("pending_approvers", tJSON, "pendingApprovers", "Approvers the current stage waits on."),
			col("current_approval_stage", tInt, "currentApprovalStage", "Current approval stage."),
			col("total_approval_stages", tInt, "totalApprovalStages", "Number of approval stages."),
			col("requested_indefinite", tBool, "requestedIndefinite", "Requested with no end date."),
			col("requested_for_seconds", tInt, "requestedForSeconds", "Requested duration."),
			col("source", tStr, "source", "member_portal, browser_extension, slack or api."),
			col("comments", tStr, "comments", "Requester comments."),
			col("requested_at", tTS, "requestedAt", "Request time."),
			col("approval_decision_made_at", tTS, "approvalDecisionMadeAt", "Decision time."),
			col("provisioned_at", tTS, "provisionedAt", "Provisioning time."),
			col("created_at", tTS, "createdAt", "Creation time."),
			col("updated_at", tTS, "updatedAt", "Last update time."),
		},
	}
}

//// ploy_integration

var integrationQuals = []qualParam{{"name", "name"}}

func tablePloyIntegration() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_integration",
		Description: "Integrations connected to Ploy (IdP, HRIS, SaaS, ...), with last scan time.",
		List:        &plugin.ListConfig{Hydrate: listHydrate("/integrations", integrationQuals), KeyColumns: optionalQuals(integrationQuals)},
		Get:         &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/integrations")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Integration ID."),
			col("name", tStr, "name", "Vendor slug, e.g. okta."),
			col("connection_name", tStr, "connectionName", "Connection label."),
			col("type", tStr, "type", "idp, email, hris, scm, infra, saas, security, knowledge_source or custom."),
			col("enabled", tBool, "enabled", "Whether enabled."),
			col("is_setup", tBool, "isSetup", "Whether setup is complete."),
			col("last_scanned", tTS, "lastScanned", "Last successful scan."),
			col("created_at", tTS, "createdAt", "Creation time."),
			col("updated_at", tTS, "updatedAt", "Last update time."),
		},
	}
}
