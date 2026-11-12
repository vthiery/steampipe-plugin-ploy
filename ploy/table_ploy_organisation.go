package ploy

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// ipWhitelistTypes are the surfaces an IP allow-list can gate (IpWhitelistType in the API spec).
var ipWhitelistTypes = []string{"dashboard", "employee-portal"}

func tablePloyIPWhitelist() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_ip_whitelist",
		Description: "IP allow-lists of the organisation, one row per surface (admin dashboard, employee portal). An empty list means the surface is unrestricted.",
		List:        &plugin.ListConfig{Hydrate: listIPWhitelists},
		Columns: []*plugin.Column{
			col("type", tStr, "type", "Surface the list gates: dashboard or employee-portal."),
			col("cidrs", tJSON, "cidrs", "Allowed ranges in CIDR notation. Empty means unrestricted."),
		},
	}
}

func listIPWhitelists(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}
	for _, t := range ipWhitelistTypes {
		row, err := client.getOne(ctx, "/organisation/ip-whitelist/"+t)
		if err != nil {
			return nil, fmt.Errorf("getting IP whitelist %s: %w", t, err)
		}
		if row == nil {
			continue
		}
		d.StreamListItem(ctx, row)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func tablePloyInternalDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_internal_domain",
		Description: "Email domains the organisation declared as internal. Members on these domains count as employees; with none declared, every member counts as internal.",
		List:        &plugin.ListConfig{Hydrate: listInternalDomains},
		Columns: []*plugin.Column{
			col("domain", tStr, "domain", "Internal email domain, without the @."),
		},
	}
}

func listInternalDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}
	row, err := client.getOne(ctx, "/organisation/internal-domains")
	if err != nil {
		return nil, fmt.Errorf("getting internal domains: %w", err)
	}
	domains, _ := row["domains"].([]interface{})
	for _, v := range domains {
		s, ok := v.(string)
		if !ok {
			continue
		}
		d.StreamListItem(ctx, Row{"domain": s})
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
