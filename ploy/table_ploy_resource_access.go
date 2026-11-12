package ploy

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"golang.org/x/sync/errgroup"
)

// fanOutConcurrency bounds parallel per-member calls on an unfiltered scan.
const fanOutConcurrency = 5

var resourceAccessQuals = []qualParam{
	{"resource_id", "resourceId"}, {"member_id", "memberId"}, {"identity_id", "identityId"}, {"active", "active"},
}

// listResourceAccess lists /resource-access. The API refuses a call with no resourceId,
// memberId or identityId; Postgres may still plan an unqualified scan (e.g. a hash join),
// so in that case fan out over every member instead of failing.
func listResourceAccess(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}
	params := paramsFromQuals(d, resourceAccessQuals)
	if params.Has("resourceId") || params.Has("memberId") || params.Has("identityId") {
		if err := client.list(ctx, d, "/resource-access", params); err != nil {
			return nil, fmt.Errorf("listing resource access: %w", err)
		}
		return nil, nil
	}

	memberIDs, err := client.allMemberIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing members for resource access fan-out: %w", err)
	}
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(fanOutConcurrency)
	for _, id := range memberIDs {
		p := url.Values{}
		for k, v := range params {
			p[k] = v
		}
		p.Set("memberId", id)
		g.Go(func() error {
			if d.RowsRemaining(gctx) == 0 {
				return nil
			}
			if err := client.list(gctx, d, "/resource-access", p); err != nil {
				return fmt.Errorf("listing resource access for member %s: %w", p.Get("memberId"), err)
			}
			return nil
		})
	}
	return nil, g.Wait()
}

// allMemberIDs pages through /members and returns every member id.
func (c *Client) allMemberIDs(ctx context.Context) ([]string, error) {
	params := url.Values{"limit": {strconv.Itoa(pageSize)}}
	var ids []string
	for offset := 0; ; {
		params.Set("offset", strconv.Itoa(offset))
		var page listResponse
		if err := c.get(ctx, "/members", params, &page); err != nil {
			return nil, err
		}
		for _, m := range page.Data {
			if id, ok := m["id"].(string); ok {
				ids = append(ids, id)
			}
		}
		offset += len(page.Data)
		if len(page.Data) == 0 || offset >= page.Pagination.Total {
			return ids, nil
		}
	}
}
