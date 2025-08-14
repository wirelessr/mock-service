package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewConfigManager tests the creation of a new ConfigManager instance
func TestNewConfigManager(t *testing.T) {
	cm := NewConfigManager()
	if cm == nil {
		t.Fatal("NewConfigManager should return a non-nil instance")
	}

	// Check that initial config is empty
	rules := cm.GetConfig()
	if len(rules) != 0 {
		t.Errorf("Expected empty rules list, got %d rules", len(rules))
	}
}

// TestLoadConfigSuccess tests successful loading of a valid JSON configuration file
func TestLoadConfigSuccess(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	configContent := `{
		"rules": [
			{
				"path": "/api/users",
				"response": {"message": "Hello Users"},
				"code": 200
			},
			{
				"path": "/api/products",
				"response": {"data": []},
				"code": 404
			}
		]
	}`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading configuration
	cm := NewConfigManager()
	err = cm.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig should succeed, got error: %v", err)
	}

	// Verify loaded configuration
	rules := cm.GetConfig()
	if len(rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(rules))
	}

	// Check first rule
	if rules[0].Path != "/api/users" {
		t.Errorf("Expected path '/api/users', got '%s'", rules[0].Path)
	}
	if rules[0].Code != 200 {
		t.Errorf("Expected code 200, got %d", rules[0].Code)
	}

	// Check second rule
	if rules[1].Path != "/api/products" {
		t.Errorf("Expected path '/api/products', got '%s'", rules[1].Path)
	}
	if rules[1].Code != 404 {
		t.Errorf("Expected code 404, got %d", rules[1].Code)
	}
}

// TestLoadConfigFileNotFound tests handling of non-existent configuration file
func TestLoadConfigFileNotFound(t *testing.T) {
	cm := NewConfigManager()
	err := cm.LoadConfig("/nonexistent/config.json")

	if err == nil {
		t.Error("LoadConfig should return error for non-existent file")
	}

	// Verify error message contains file path
	if err != nil && len(err.Error()) == 0 {
		t.Error("Error message should not be empty")
	}
}

// TestLoadConfigInvalidJSON tests handling of invalid JSON format
func TestLoadConfigInvalidJSON(t *testing.T) {
	// Create temporary file with invalid JSON
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid.json")

	invalidJSON := `{
		"rules": [
			{
				"path": "/api/users"
				"response": {"message": "Hello Users"}  // Missing comma
			}
		]
	}`

	err := os.WriteFile(configFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading invalid configuration
	cm := NewConfigManager()
	err = cm.LoadConfig(configFile)

	if err == nil {
		t.Error("LoadConfig should return error for invalid JSON")
	}
}

// TestGetConfigAfterLoad tests that GetConfig returns the correct rules after loading
func TestGetConfigAfterLoad(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	configContent := `{
		"rules": [
			{
				"path": "/test",
				"response": {"status": "ok"},
				"code": 201
			}
		]
	}`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cm := NewConfigManager()

	// Before loading, should be empty
	rules := cm.GetConfig()
	if len(rules) != 0 {
		t.Errorf("Expected 0 rules before loading, got %d", len(rules))
	}

	// After loading, should contain the rule
	err = cm.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	rules = cm.GetConfig()
	if len(rules) != 1 {
		t.Errorf("Expected 1 rule after loading, got %d", len(rules))
	}

	if rules[0].Path != "/test" {
		t.Errorf("Expected path '/test', got '%s'", rules[0].Path)
	}
}
