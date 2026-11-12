package ploy

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePloyAuditLog() *plugin.Table {
	return &plugin.Table{
		Name:        "ploy_audit_log",
		Description: "The organisation's Ploy audit trail. Filter occurred_at (>, >=, <, <=) to bound the scan.",
		List: &plugin.ListConfig{
			Hydrate: listAuditLogs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "type", Require: plugin.Optional},
				{Name: "actor_user_id", Require: plugin.Optional},
				{Name: "occurred_at", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
			},
		},
		Get: &plugin.GetConfig{KeyColumns: plugin.SingleColumn("id"), Hydrate: getHydrate("/audit-logs")},
		Columns: []*plugin.Column{
			col("id", tStr, "id", "Entry ID."),
			col("type", tStr, "type", "Event type, e.g. workflows.executed.create_resource."),
			col("occurred_at", tTS, "occurredAt", "When the action started."),
			col("action_result", tStr, "actionResult", "success or failed; null for events with no outcome."),
			col("actor_type", tStr, "actor.type", "user, member, api_key or system."),
			col("actor_user_id", tStr, "actor.userId", "Acting user ID."),
			col("actor_member_id", tStr, "actor.memberId", "Acting member ID."),
			col("actor_api_key_id", tStr, "actor.apiKeyId", "Acting API key ID."),
			col("actor_email", tStr, "actor.email", "Acting user's/member's email."),
			col("actor_name", tStr, "actor.name", "Acting user's/member's name (`Ploy Staff` for Ploy support)."),
			col("actor", tJSON, "actor", "Actor object."),
			col("ip_address", tStr, "ipAddress", "Source IP."),
			col("user_agent", tStr, "userAgent", "Source user agent."),
			col("payload", tJSON, "payload", "Event details (secrets redacted)."),
		},
	}
}

func listAuditLogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}
	params := paramsFromQuals(d, []qualParam{{"type", "type"}, {"actor_user_id", "userId"}})

	// occurredSince is inclusive, occurredBefore exclusive, both whole seconds.
	// Widen to whole seconds; Postgres re-applies the exact predicate.
	if q := d.Quals["occurred_at"]; q != nil {
		for _, qual := range q.Quals {
			ts := qual.Value.GetTimestampValue()
			if ts == nil {
				continue
			}
			t := ts.AsTime().UTC()
			switch qual.Operator {
			case ">", ">=":
				params.Set("occurredSince", t.Truncate(time.Second).Format(time.RFC3339))
			case "<", "<=":
				params.Set("occurredBefore", t.Truncate(time.Second).Add(time.Second).Format(time.RFC3339))
			case "=":
				params.Set("occurredSince", t.Truncate(time.Second).Format(time.RFC3339))
				params.Set("occurredBefore", t.Truncate(time.Second).Add(time.Second).Format(time.RFC3339))
			}
		}
	}

	if err := client.list(ctx, d, "/audit-logs", params); err != nil {
		return nil, fmt.Errorf("listing audit logs: %w", err)
	}
	return nil, nil
}
