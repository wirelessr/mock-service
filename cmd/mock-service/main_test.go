package main

import (
	"mock-service/internal/config"
	"mock-service/internal/logger"
	"mock-service/internal/matcher"
	"mock-service/internal/response"
	"os"
	"path/filepath"
	"testing"
)

// TestMainProgramIntegration tests the complete integration of the main program
func TestMainProgramIntegration(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test_config.json")

	configContent := `{
		"rules": [
			{
				"path": "/api/users",
				"response": {
					"users": ["alice", "bob", "charlie"],
					"count": 3
				},
				"code": 200
			},
			{
				"path": "/api/products",
				"response": {
					"products": [],
					"message": "No products found"
				},
				"code": 404
			}
		]
	}`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test configuration loading (without starting server)
	// This tests the main components integration
	testConfigLoading(t, configFile)
}

// testConfigLoading tests that configuration can be loaded successfully
func testConfigLoading(t *testing.T, configFile string) {
	// This simulates the main program's configuration loading logic
	configManager := config.NewConfigManager()

	err := configManager.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Configuration loading should succeed, got error: %v", err)
	}

	rules := configManager.GetConfig()
	if len(rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(rules))
	}

	// Verify first rule
	if rules[0].Path != "/api/users" {
		t.Errorf("Expected first rule path '/api/users', got '%s'", rules[0].Path)
	}

	if rules[0].Code != 200 {
		t.Errorf("Expected first rule code 200, got %d", rules[0].Code)
	}

	// Verify second rule
	if rules[1].Path != "/api/products" {
		t.Errorf("Expected second rule path '/api/products', got '%s'", rules[1].Path)
	}

	if rules[1].Code != 404 {
		t.Errorf("Expected second rule code 404, got %d", rules[1].Code)
	}
}

// TestMainProgramWithInvalidConfig tests handling of invalid configuration
func TestMainProgramWithInvalidConfig(t *testing.T) {
	// Create temporary invalid config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid_config.json")

	invalidConfigContent := `{
		"rules": [
			{
				"path": "/api/test"
				"response": {"message": "test"}  // Missing comma
			}
		]
	}`

	err := os.WriteFile(configFile, []byte(invalidConfigContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid config file: %v", err)
	}

	// Test that invalid configuration is handled properly
	configManager := config.NewConfigManager()

	err = configManager.LoadConfig(configFile)
	if err == nil {
		t.Error("Expected error when loading invalid configuration")
	}
}

// TestMainProgramWithNonExistentConfig tests handling of non-existent configuration file
func TestMainProgramWithNonExistentConfig(t *testing.T) {
	configManager := config.NewConfigManager()

	err := configManager.LoadConfig("/nonexistent/config.json")
	if err == nil {
		t.Error("Expected error when loading non-existent configuration file")
	}
}

// TestComponentsIntegration tests that all components work together correctly
func TestComponentsIntegration(t *testing.T) {
	// Create test configuration
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "integration_config.json")

	configContent := `{
		"rules": [
			{
				"path": "/test/integration",
				"response": {
					"status": "success",
					"data": {"id": 123, "name": "integration test"}
				},
				"code": 201
			}
		]
	}`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create integration config file: %v", err)
	}

	// Initialize all components as in main
	configManager := config.NewConfigManager()
	pathMatcher := matcher.NewPathMatcher()
	responseBuilder := response.NewResponseBuilder()
	_ = logger.NewLogger() // Logger not used in this test but part of integration

	// Load configuration
	err = configManager.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test path matching
	rules := configManager.GetConfig()
	rule, found := pathMatcher.FindMatch("/test/integration", rules)
	if !found {
		t.Fatal("Expected to find matching rule")
	}

	// Test response building
	statusCode, body := responseBuilder.BuildResponse(rule)
	if statusCode != 201 {
		t.Errorf("Expected status code 201, got %d", statusCode)
	}

	// Verify response structure
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		t.Fatal("Expected response body to be a map")
	}

	if bodyMap["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", bodyMap["status"])
	}

	// Test default response
	_, found = pathMatcher.FindMatch("/nonexistent", rules)
	if found {
		t.Error("Expected no match for non-existent path")
	}

	defaultStatusCode, defaultBody := responseBuilder.BuildDefaultResponse()
	if defaultStatusCode != 200 {
		t.Errorf("Expected default status code 200, got %d", defaultStatusCode)
	}

	defaultBodyMap, ok := defaultBody.(map[string]interface{})
	if !ok {
		t.Fatal("Expected default response body to be a map")
	}

	// Should be empty JSON object
	if len(defaultBodyMap) != 0 {
		t.Errorf("Expected empty JSON object, got %v", defaultBodyMap)
	}
}
