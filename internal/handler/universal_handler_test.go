package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mock-service/internal/models"

	"github.com/gin-gonic/gin"
)

// Mock implementations for testing
type mockConfigManager struct {
	rules []models.MockRule
}

func (m *mockConfigManager) LoadConfig(filePath string) error {
	return nil
}

func (m *mockConfigManager) GetConfig() []models.MockRule {
	return m.rules
}

type mockPathMatcher struct {
	shouldMatch  bool
	ruleToReturn *models.MockRule
}

func (m *mockPathMatcher) FindMatch(requestPath string, rules []models.MockRule) (*models.MockRule, bool) {
	return m.ruleToReturn, m.shouldMatch
}

type mockResponseBuilder struct{}

func (m *mockResponseBuilder) BuildResponse(rule *models.MockRule) (statusCode int, body interface{}) {
	return rule.Code, rule.Response
}

func (m *mockResponseBuilder) BuildDefaultResponse() (statusCode int, body interface{}) {
	return 200, map[string]interface{}{}
}

type mockLogger struct {
	loggedRequests  []LoggedRequest
	loggedResponses []LoggedResponse
	loggedMatches   []*models.MockRule
	defaultLogged   bool
}

type LoggedRequest struct {
	Method string
	Path   string
	Params map[string]string
}

type LoggedResponse struct {
	StatusCode int
	Body       interface{}
}

func (m *mockLogger) LogRequest(method, path string, params map[string]string) {
	m.loggedRequests = append(m.loggedRequests, LoggedRequest{
		Method: method,
		Path:   path,
		Params: params,
	})
}

func (m *mockLogger) LogResponse(statusCode int, body interface{}) {
	m.loggedResponses = append(m.loggedResponses, LoggedResponse{
		StatusCode: statusCode,
		Body:       body,
	})
}

func (m *mockLogger) LogMatch(rule *models.MockRule) {
	m.loggedMatches = append(m.loggedMatches, rule)
}

func (m *mockLogger) LogDefault() {
	m.defaultLogged = true
}

// TestNewUniversalHandler tests the creation of a new UniversalHandler instance
func TestNewUniversalHandler(t *testing.T) {
	configManager := &mockConfigManager{}
	pathMatcher := &mockPathMatcher{}
	responseBuilder := &mockResponseBuilder{}
	logger := &mockLogger{}

	handler := NewUniversalHandler(configManager, pathMatcher, responseBuilder, logger)

	if handler == nil {
		t.Fatal("NewUniversalHandler should return a non-nil instance")
	}
}

// TestHandleRequestWithMatchingRule tests handling request when a rule matches
func TestHandleRequestWithMatchingRule(t *testing.T) {
	// Set up mocks
	rule := &models.MockRule{
		Path: "/api/users",
		Response: map[string]interface{}{
			"users": []string{"alice", "bob"},
		},
		Code: 200,
	}

	configManager := &mockConfigManager{
		rules: []models.MockRule{*rule},
	}
	pathMatcher := &mockPathMatcher{
		shouldMatch:  true,
		ruleToReturn: rule,
	}
	responseBuilder := &mockResponseBuilder{}
	logger := &mockLogger{}

	handler := NewUniversalHandler(configManager, pathMatcher, responseBuilder, logger)

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Any("/*path", handler.HandleRequest)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/users?id=123&name=test", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// Check response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	users, ok := responseBody["users"].([]interface{})
	if !ok {
		t.Fatal("Expected users array in response")
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Check logging
	if len(logger.loggedRequests) != 1 {
		t.Errorf("Expected 1 logged request, got %d", len(logger.loggedRequests))
	}

	if logger.loggedRequests[0].Method != "GET" {
		t.Errorf("Expected method GET, got %s", logger.loggedRequests[0].Method)
	}

	if logger.loggedRequests[0].Path != "/api/users" {
		t.Errorf("Expected path /api/users, got %s", logger.loggedRequests[0].Path)
	}

	// Check query parameters were logged
	if logger.loggedRequests[0].Params["id"] != "123" {
		t.Errorf("Expected param id=123, got %s", logger.loggedRequests[0].Params["id"])
	}

	if len(logger.loggedMatches) != 1 {
		t.Errorf("Expected 1 logged match, got %d", len(logger.loggedMatches))
	}

	if len(logger.loggedResponses) != 1 {
		t.Errorf("Expected 1 logged response, got %d", len(logger.loggedResponses))
	}

	if logger.defaultLogged {
		t.Error("Expected default not to be logged when rule matches")
	}
}

// TestHandleRequestWithNoMatchingRule tests handling request when no rule matches
func TestHandleRequestWithNoMatchingRule(t *testing.T) {
	// Set up mocks
	configManager := &mockConfigManager{
		rules: []models.MockRule{},
	}
	pathMatcher := &mockPathMatcher{
		shouldMatch:  false,
		ruleToReturn: nil,
	}
	responseBuilder := &mockResponseBuilder{}
	logger := &mockLogger{}

	handler := NewUniversalHandler(configManager, pathMatcher, responseBuilder, logger)

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Any("/*path", handler.HandleRequest)

	// Create test request
	req, _ := http.NewRequest("POST", "/api/unknown", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Check response status should be 200 as per requirements
	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// Check response body should be empty JSON object
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Should be empty JSON object
	if len(responseBody) != 0 {
		t.Errorf("Expected empty JSON object, got %v", responseBody)
	}

	// Check logging
	if len(logger.loggedRequests) != 1 {
		t.Errorf("Expected 1 logged request, got %d", len(logger.loggedRequests))
	}

	if logger.loggedRequests[0].Method != "POST" {
		t.Errorf("Expected method POST, got %s", logger.loggedRequests[0].Method)
	}

	if len(logger.loggedMatches) != 0 {
		t.Errorf("Expected 0 logged matches, got %d", len(logger.loggedMatches))
	}

	if !logger.defaultLogged {
		t.Error("Expected default to be logged when no rule matches")
	}

	if len(logger.loggedResponses) != 1 {
		t.Errorf("Expected 1 logged response, got %d", len(logger.loggedResponses))
	}
}

// TestHandleRequestWithQueryParameters tests query parameter parsing
func TestHandleRequestWithQueryParameters(t *testing.T) {
	// Set up mocks
	configManager := &mockConfigManager{}
	pathMatcher := &mockPathMatcher{shouldMatch: false}
	responseBuilder := &mockResponseBuilder{}
	logger := &mockLogger{}

	handler := NewUniversalHandler(configManager, pathMatcher, responseBuilder, logger)

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Any("/*path", handler.HandleRequest)

	// Create test request with multiple query parameters
	req, _ := http.NewRequest("GET", "/api/test?param1=value1&param2=value2&param3=value3", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Check that query parameters were parsed correctly
	if len(logger.loggedRequests) != 1 {
		t.Fatalf("Expected 1 logged request, got %d", len(logger.loggedRequests))
	}

	params := logger.loggedRequests[0].Params

	if params["param1"] != "value1" {
		t.Errorf("Expected param1=value1, got %s", params["param1"])
	}

	if params["param2"] != "value2" {
		t.Errorf("Expected param2=value2, got %s", params["param2"])
	}

	if params["param3"] != "value3" {
		t.Errorf("Expected param3=value3, got %s", params["param3"])
	}
}

// TestHandleRequestWithMultipleQueryValues tests handling of multiple values for same parameter
func TestHandleRequestWithMultipleQueryValues(t *testing.T) {
	// Set up mocks
	configManager := &mockConfigManager{}
	pathMatcher := &mockPathMatcher{shouldMatch: false}
	responseBuilder := &mockResponseBuilder{}
	logger := &mockLogger{}

	handler := NewUniversalHandler(configManager, pathMatcher, responseBuilder, logger)

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Any("/*path", handler.HandleRequest)

	// Create test request with multiple values for same parameter
	req, _ := http.NewRequest("GET", "/api/test?tags=tag1&tags=tag2&tags=tag3", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Check that only first value is taken for parameters with multiple values
	if len(logger.loggedRequests) != 1 {
		t.Fatalf("Expected 1 logged request, got %d", len(logger.loggedRequests))
	}

	params := logger.loggedRequests[0].Params

	if params["tags"] != "tag1" {
		t.Errorf("Expected tags=tag1 (first value), got %s", params["tags"])
	}
}
