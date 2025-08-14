package response

import (
	"reflect"
	"testing"

	"mock-service/internal/models"
)

// TestNewResponseBuilder tests the creation of a new ResponseBuilder instance
func TestNewResponseBuilder(t *testing.T) {
	rb := NewResponseBuilder()
	if rb == nil {
		t.Fatal("NewResponseBuilder should return a non-nil instance")
	}
}

// TestBuildResponseWithValidRule tests building response from a valid mock rule
func TestBuildResponseWithValidRule(t *testing.T) {
	rb := NewResponseBuilder()

	rule := &models.MockRule{
		Path: "/api/users",
		Response: map[string]interface{}{
			"users": []string{"alice", "bob"},
			"count": 2,
		},
		Code: 200,
	}

	statusCode, body := rb.BuildResponse(rule)

	// Check status code
	if statusCode != 200 {
		t.Errorf("Expected status code 200, got %d", statusCode)
	}

	// Check response body
	expectedBody := map[string]interface{}{
		"users": []string{"alice", "bob"},
		"count": 2,
	}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Expected body %v, got %v", expectedBody, body)
	}
}

// TestBuildResponseWithCustomStatusCode tests building response with custom status code
func TestBuildResponseWithCustomStatusCode(t *testing.T) {
	rb := NewResponseBuilder()

	rule := &models.MockRule{
		Path: "/api/error",
		Response: map[string]interface{}{
			"error": "Internal Server Error",
		},
		Code: 500,
	}

	statusCode, body := rb.BuildResponse(rule)

	// Check status code
	if statusCode != 500 {
		t.Errorf("Expected status code 500, got %d", statusCode)
	}

	// Check response body
	expectedBody := map[string]interface{}{
		"error": "Internal Server Error",
	}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Expected body %v, got %v", expectedBody, body)
	}
}

// TestBuildResponseWithZeroStatusCode tests default status code when rule has zero code
func TestBuildResponseWithZeroStatusCode(t *testing.T) {
	rb := NewResponseBuilder()

	rule := &models.MockRule{
		Path: "/api/test",
		Response: map[string]interface{}{
			"message": "test response",
		},
		Code: 0, // Zero status code should default to 200
	}

	statusCode, body := rb.BuildResponse(rule)

	// Should default to 200
	if statusCode != 200 {
		t.Errorf("Expected default status code 200, got %d", statusCode)
	}

	// Check response body
	expectedBody := map[string]interface{}{
		"message": "test response",
	}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Expected body %v, got %v", expectedBody, body)
	}
}

// TestBuildResponseWithEmptyResponse tests building response with empty response body
func TestBuildResponseWithEmptyResponse(t *testing.T) {
	rb := NewResponseBuilder()

	rule := &models.MockRule{
		Path:     "/api/empty",
		Response: map[string]interface{}{},
		Code:     204,
	}

	statusCode, body := rb.BuildResponse(rule)

	// Check status code
	if statusCode != 204 {
		t.Errorf("Expected status code 204, got %d", statusCode)
	}

	// Check response body is empty map
	expectedBody := map[string]interface{}{}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Expected empty body %v, got %v", expectedBody, body)
	}
}

// TestBuildDefaultResponse tests building default response when no rule matches
func TestBuildDefaultResponse(t *testing.T) {
	rb := NewResponseBuilder()

	statusCode, body := rb.BuildDefaultResponse()

	// Check status code should be 200 as per requirements
	if statusCode != 200 {
		t.Errorf("Expected default status code 200, got %d", statusCode)
	}

	// Check response body structure - should be empty JSON object
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		t.Fatal("Expected response body to be a map")
	}

	// Should be empty JSON object
	if len(bodyMap) != 0 {
		t.Errorf("Expected empty JSON object, got %v", bodyMap)
	}
}

// TestBuildResponseWithComplexData tests building response with complex nested data
func TestBuildResponseWithComplexData(t *testing.T) {
	rb := NewResponseBuilder()

	rule := &models.MockRule{
		Path: "/api/complex",
		Response: map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"id":   123,
					"name": "John Doe",
					"tags": []string{"admin", "user"},
				},
				"metadata": map[string]interface{}{
					"version": "1.0",
					"active":  true,
				},
			},
			"status": "success",
		},
		Code: 201,
	}

	statusCode, body := rb.BuildResponse(rule)

	// Check status code
	if statusCode != 201 {
		t.Errorf("Expected status code 201, got %d", statusCode)
	}

	// Check that complex data structure is preserved
	if !reflect.DeepEqual(body, rule.Response) {
		t.Errorf("Expected body to match rule response exactly")
	}
}
