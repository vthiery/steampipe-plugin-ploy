connection "ploy" {
  plugin = "ghcr.io/vthiery/ploy"

  # API key for the Ploy public API. Read-only scopes are sufficient.
  # Create one in the Ploy dashboard. Falls back to the PLOY_API_KEY
  # environment variable when not set here.
  # api_key = "YOUR_API_KEY"

  # Base URL of the API. Defaults to https://api.joinploy.com/api/v1.
  # base_url = "https://api.joinploy.com/api/v1"
}
