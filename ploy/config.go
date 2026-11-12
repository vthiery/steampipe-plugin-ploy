package ploy

import "github.com/turbot/steampipe-plugin-sdk/v5/plugin"

// PloyConfig holds the connection configuration for the Ploy plugin.
type PloyConfig struct {
	// API key for the Ploy public API (Bearer token). Read-only scopes are enough.
	APIKey *string `hcl:"api_key"`
	// Base URL of the API. Defaults to https://api.joinploy.com/api/v1.
	BaseURL *string `hcl:"base_url"`
}

// ConfigInstance returns a new, empty PloyConfig.
func ConfigInstance() interface{} {
	return &PloyConfig{}
}

// GetConfig retrieves the connection configuration for the given connection.
func GetConfig(connection *plugin.Connection) PloyConfig {
	if connection == nil || connection.Config == nil {
		return PloyConfig{}
	}
	cfg, _ := connection.Config.(PloyConfig)
	return cfg
}
