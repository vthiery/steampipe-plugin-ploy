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

// optionalQuals turns qual mappings into optional key columns.
func optionalQuals(qs []qualParam) []*plugin.KeyColumn {
	kcs := make([]*plugin.KeyColumn, 0, len(qs))
	for _, q := range qs {
		kcs = append(kcs, &plugin.KeyColumn{Name: q.Column, Require: plugin.Optional})
	}
	return kcs
}

// paramsFromQuals copies equality quals into API query parameters.
func paramsFromQuals(d *plugin.QueryData, qs []qualParam) url.Values {
	params := url.Values{}
	for _, q := range qs {
		v, ok := d.EqualsQuals[q.Column]
		if !ok || v == nil {
			continue
		}
		switch x := v.Value.(type) {
		case *proto.QualValue_StringValue:
			params.Set(q.Param, x.StringValue)
		case *proto.QualValue_BoolValue:
			params.Set(q.Param, strconv.FormatBool(x.BoolValue))
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
