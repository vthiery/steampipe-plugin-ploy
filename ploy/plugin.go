// Package ploy provides a Steampipe plugin for the Ploy identity governance public API.
package ploy

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the definition of the Ploy Steampipe plugin.
func Plugin(ctx context.Context) *plugin.Plugin {
	return &plugin.Plugin{
		Name:             "steampipe-plugin-ploy",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"ploy_member":                 tablePloyMember(),
			"ploy_user":                   tablePloyUser(),
			"ploy_resource":               tablePloyResource(),
			"ploy_resource_access":        tablePloyResourceAccess(),
			"ploy_access_review_campaign": tablePloyAccessReviewCampaign(),
			"ploy_access_review":          tablePloyAccessReview(),
			"ploy_access_request":         tablePloyAccessRequest(),
			"ploy_audit_log":              tablePloyAuditLog(),
			"ploy_integration":            tablePloyIntegration(),
			"ploy_ip_whitelist":           tablePloyIPWhitelist(),
			"ploy_internal_domain":        tablePloyInternalDomain(),
		},
	}
}
