package handler

import (
	"mock-service/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// UniversalHandler handles all HTTP requests using the configured components
type UniversalHandler struct {
	configManager   interfaces.ConfigManager
	pathMatcher     interfaces.PathMatcher
	responseBuilder interfaces.ResponseBuilder
	logger          interfaces.Logger
}

// NewUniversalHandler creates a new instance of UniversalHandler
func NewUniversalHandler(
	configManager interfaces.ConfigManager,
	pathMatcher interfaces.PathMatcher,
	responseBuilder interfaces.ResponseBuilder,
	logger interfaces.Logger,
) *UniversalHandler {
	return &UniversalHandler{
		configManager:   configManager,
		pathMatcher:     pathMatcher,
		responseBuilder: responseBuilder,
		logger:          logger,
	}
}

// HandleRequest handles all HTTP requests for any path and method
func (uh *UniversalHandler) HandleRequest(c *gin.Context) {
	// Extract request information
	method := c.Request.Method
	path := c.Request.URL.Path

	// Parse query parameters
	params := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0] // Take first value if multiple exist
		}
	}

	// Log the incoming request
	uh.logger.LogRequest(method, path, params)

	// Get current configuration rules
	rules := uh.configManager.GetConfig()

	// Try to find a matching rule
	rule, found := uh.pathMatcher.FindMatch(path, rules)

	var statusCode int
	var body interface{}

	if found {
		// Rule matched - build response from rule
		uh.logger.LogMatch(rule)
		statusCode, body = uh.responseBuilder.BuildResponse(rule)
	} else {
		// No rule matched - use default response
		uh.logger.LogDefault()
		statusCode, body = uh.responseBuilder.BuildDefaultResponse()
	}

	// Log the response
	uh.logger.LogResponse(statusCode, body)

	// Send the response
	c.JSON(statusCode, body)
}
