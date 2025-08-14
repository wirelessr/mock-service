package models

// MockRule represents a single mock rule configuration
// It defines how the service should respond to requests matching a specific path
type MockRule struct {
	// Path is the request path to match against (e.g., "/api/users")
	Path string `json:"path"`
	// Response is the JSON response body to return when this rule matches
	Response map[string]interface{} `json:"response"`
	// Code is the HTTP status code to return (defaults to 200 if not specified)
	Code int `json:"code"`
}

// Config represents the complete configuration structure loaded from JSON file
// It contains all the mock rules that define the service behavior
type Config struct {
	// Rules is the list of mock rules to be processed in order
	Rules []MockRule `json:"rules"`
}
