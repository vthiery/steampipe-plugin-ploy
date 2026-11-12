package ploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"golang.org/x/time/rate"
)

const (
	defaultBaseURL = "https://api.joinploy.com/api/v1"
	pageSize       = 100 // API maximum
	maxRetries     = 8
)

// limiter is shared by every request the plugin process makes. Ploy allows 600 requests
// per API key per 60-second window (X-RateLimit-Limit); 8 req/s stays under that and
// leaves headroom for other clients using the same key.
var limiter = rate.NewLimiter(rate.Limit(8), 8)

// Client wraps the Ploy public REST API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Row is a single API object. Columns read fields with transform.FromField("camelCaseKey").
type Row = map[string]interface{}

type pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type listResponse struct {
	Data       []Row      `json:"data"`
	Pagination pagination `json:"pagination"`
}

type getResponse struct {
	Data Row `json:"data"`
}

func getClient(_ context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)
	apiKey := os.Getenv("PLOY_API_KEY")
	if cfg.APIKey != nil && *cfg.APIKey != "" {
		apiKey = *cfg.APIKey
	}
	if apiKey == "" {
		return nil, fmt.Errorf("api_key must be configured (or set PLOY_API_KEY)")
	}
	base := defaultBaseURL
	if cfg.BaseURL != nil && *cfg.BaseURL != "" {
		base = strings.TrimRight(*cfg.BaseURL, "/")
	}
	return &Client{apiKey: apiKey, baseURL: base, httpClient: &http.Client{Timeout: 60 * time.Second}}, nil
}

// get performs an authenticated GET, retrying on 429 and 5xx, and decodes the body into result.
func (c *Client) get(ctx context.Context, path string, params url.Values, result interface{}) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := limiter.Wait(ctx); err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return fmt.Errorf("creating request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("executing request: %w", err)
		} else {
			body, readErr := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			switch {
			case readErr != nil:
				lastErr = fmt.Errorf("reading response body: %w", readErr)
			case resp.StatusCode == http.StatusOK:
				if err := json.Unmarshal(body, result); err != nil {
					return fmt.Errorf("unmarshalling response from %s: %w", path, err)
				}
				return nil
			case resp.StatusCode == http.StatusNotFound:
				return errNotFound
			case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500:
				lastErr = fmt.Errorf("GET %s: status %d: %s", path, resp.StatusCode, truncate(body))
				if wait := retryAfter(resp); wait > 0 {
					if err := sleep(ctx, wait); err != nil {
						return err
					}
					continue
				}
			default:
				return fmt.Errorf("GET %s: status %d: %s", path, resp.StatusCode, truncate(body))
			}
		}
		// 0.5s, 1s, 2s, ... capped at 30s.
		backoff := time.Duration(1<<attempt) * 500 * time.Millisecond
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
		if err := sleep(ctx, backoff); err != nil {
			return err
		}
	}
	return lastErr
}

var errNotFound = fmt.Errorf("not found")

// list pages through an offset-paginated endpoint, streaming every row.
func (c *Client) list(ctx context.Context, d *plugin.QueryData, path string, params url.Values) error {
	filters := filterValues(params)
	if params == nil {
		params = url.Values{}
	}
	params.Set("limit", strconv.Itoa(pageSize))
	offset := 0
	for {
		params.Set("offset", strconv.Itoa(offset))
		var page listResponse
		if err := c.get(ctx, path, params, &page); err != nil {
			return err
		}
		for _, row := range page.Data {
			for k, v := range filters {
				if _, ok := row[k]; !ok {
					row[k] = v
				}
			}
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
		offset += len(page.Data)
		if len(page.Data) == 0 || offset >= page.Pagination.Total {
			return nil
		}
	}
}

// getOne fetches a single object by path, returning nil (no row) on 404.
func (c *Client) getOne(ctx context.Context, path string) (Row, error) {
	var res getResponse
	if err := c.get(ctx, path, nil, &res); err != nil {
		if err == errNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res.Data, nil
}

// retryAfter returns how long Ploy asks us to wait: Retry-After, or X-RateLimit-Reset
// (both in seconds until the current window resets).
func retryAfter(resp *http.Response) time.Duration {
	for _, h := range []string{"Retry-After", "X-RateLimit-Reset"} {
		if s, err := strconv.Atoi(resp.Header.Get(h)); err == nil && s > 0 {
			return time.Duration(s) * time.Second
		}
	}
	return 0
}

func sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

func truncate(b []byte) string {
	if len(b) > 300 {
		return string(b[:300]) + "…"
	}
	return string(b)
}

// filterValues returns the request filters as typed values keyed by API parameter name.
// list copies them onto rows that do not already carry the field, so filter-only columns
// (e.g. memberStatus on /resource-access) report the value the query filtered on.
func filterValues(params url.Values) Row {
	out := Row{}
	for k, v := range params {
		if len(v) == 0 || k == "limit" || k == "offset" {
			continue
		}
		switch v[0] {
		case "true":
			out[k] = true
		case "false":
			out[k] = false
		default:
			out[k] = v[0]
		}
	}
	return out
}
