package ploy

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// col builds a column that reads the camelCase API field `field`.
func col(name string, t proto.ColumnType, field, desc string) *plugin.Column {
	return &plugin.Column{Name: name, Type: t, Transform: transform.FromField(field), Description: desc}
}

// qualParam maps an equality qual on a column to an API query parameter.
type qualParam struct {
	Column string
	Param  string
}

// filterOnly lists columns the API does not return: they only report the value the query
// filtered on. They need CacheMatch "exact", otherwise the query cache answers a filtered
// query from cached rows fetched without the filter, where the column is null, and the
// filter then removes every row.
var filterOnly = map[string]bool{
	"member_status": true, "member_is_managed": true, "has_member": true,
	"identity_type": true, "resource_type": true, "managed_app": true,
}

// optionalQuals turns qual mappings into optional key columns. Boolean columns also accept
// "<>": Postgres rewrites `col = false` into `col <> true`, and without it the filter would
// not be pushed down and the API would be scanned in full.
func optionalQuals(qs []qualParam, boolCols ...string) []*plugin.KeyColumn {
	isBool := map[string]bool{}
	for _, c := range boolCols {
		isBool[c] = true
	}
	kcs := make([]*plugin.KeyColumn, 0, len(qs))
	for _, q := range qs {
		kc := &plugin.KeyColumn{Name: q.Column, Require: plugin.Optional}
		if filterOnly[q.Column] {
			kc.CacheMatch = "exact"
		}
		if isBool[q.Column] {
			kc.Operators = []string{"=", "<>"}
		}
		kcs = append(kcs, kc)
	}
	return kcs
}

// paramsFromQuals copies equality quals into API query parameters.
func paramsFromQuals(d *plugin.QueryData, qs []qualParam) url.Values {
	params := url.Values{}
	for _, q := range qs {
		quals := d.Quals[q.Column]
		if quals == nil {
			continue
		}
		for _, qual := range quals.Quals {
			if qual.Value == nil {
				continue
			}
			switch x := qual.Value.Value.(type) {
			case *proto.QualValue_StringValue:
				if qual.Operator == "=" {
					params.Set(q.Param, x.StringValue)
				}
			case *proto.QualValue_BoolValue:
				switch qual.Operator {
				case "=":
					params.Set(q.Param, strconv.FormatBool(x.BoolValue))
				case "<>":
					params.Set(q.Param, strconv.FormatBool(!x.BoolValue))
				}
			}
		}
	}
	return params
}

// listHydrate returns a List hydrate that pages through path, pushing down the given quals.
func listHydrate(path string, qs []qualParam) plugin.HydrateFunc {
	return func(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
		client, err := getClient(ctx, d)
		if err != nil {
			return nil, err
		}
		if err := client.list(ctx, d, path, paramsFromQuals(d, qs)); err != nil {
			return nil, fmt.Errorf("listing %s: %w", path, err)
		}
		return nil, nil
	}
}

// getHydrate returns a Get hydrate fetching basePath/{id}.
func getHydrate(basePath string) plugin.HydrateFunc {
	return func(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
		id := d.EqualsQualString("id")
		if id == "" {
			return nil, nil
		}
		client, err := getClient(ctx, d)
		if err != nil {
			return nil, err
		}
		row, err := client.getOne(ctx, basePath+"/"+url.PathEscape(id))
		if err != nil {
			return nil, fmt.Errorf("getting %s/%s: %w", basePath, id, err)
		}
		if row == nil {
			return nil, nil
		}
		return row, nil
	}
}
