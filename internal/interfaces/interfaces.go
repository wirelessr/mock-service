package interfaces

import "mock-service/internal/models"

// ConfigManager handles loading and managing JSON configuration files
type ConfigManager interface {
	// LoadConfig loads configuration from the specified file path
	LoadConfig(filePath string) error
	// GetConfig returns the current list of mock rules
	GetConfig() []models.MockRule
}

// PathMatcher handles matching request paths against configured rules
type PathMatcher interface {
	// FindMatch finds the first matching rule for the given request path
	// Returns the matched rule and true if found, nil and false otherwise
	FindMatch(requestPath string, rules []models.MockRule) (*models.MockRule, bool)
}

// ResponseBuilder handles building HTTP responses based on mock rules
type ResponseBuilder interface {
	// BuildResponse builds a response based on the provided mock rule
	BuildResponse(rule *models.MockRule) (statusCode int, body interface{})
	// BuildDefaultResponse builds a default response when no rule matches
	BuildDefaultResponse() (statusCode int, body interface{})
}

// Logger provides structured logging functionality for the mock service
type Logger interface {
	// LogRequest logs incoming HTTP request details
	LogRequest(method, path string, params map[string]string)
	// LogResponse logs outgoing HTTP response details
	LogResponse(statusCode int, body interface{})
	// LogMatch logs when a rule is matched
	LogMatch(rule *models.MockRule)
	// LogDefault logs when default response is used
	LogDefault()
}
