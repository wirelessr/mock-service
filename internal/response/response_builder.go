package response

import "mock-service/internal/models"

// ResponseBuilderImpl implements the ResponseBuilder interface
// It handles building HTTP responses based on mock rules
type ResponseBuilderImpl struct{}

// NewResponseBuilder creates a new instance of ResponseBuilder
func NewResponseBuilder() *ResponseBuilderImpl {
	return &ResponseBuilderImpl{}
}

// BuildResponse builds a response based on the provided mock rule
// Returns the status code and response body from the rule
func (rb *ResponseBuilderImpl) BuildResponse(rule *models.MockRule) (statusCode int, body interface{}) {
	// Use the status code from the rule, default to 200 if not specified or invalid
	statusCode = rule.Code
	if statusCode == 0 {
		statusCode = 200
	}

	// Return the response body from the rule
	body = rule.Response
	return statusCode, body
}

// BuildDefaultResponse builds a default response when no rule matches
// Returns 200 status with an empty JSON object as per requirements
func (rb *ResponseBuilderImpl) BuildDefaultResponse() (statusCode int, body interface{}) {
	statusCode = 200
	body = map[string]interface{}{}
	return statusCode, body
}
