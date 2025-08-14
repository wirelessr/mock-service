package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"mock-service/internal/models"
)

// captureOutput captures stdout output during test execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestNewLogger tests the creation of a new Logger instance
func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("NewLogger should return a non-nil instance")
	}
}

// TestLogRequest tests logging of HTTP request details
func TestLogRequest(t *testing.T) {
	logger := NewLogger()

	params := map[string]string{
		"id":   "123",
		"name": "test",
	}

	output := captureOutput(func() {
		logger.LogRequest("GET", "/api/users", params)
	})

	// Verify output is valid JSON
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
	if err != nil {
		t.Fatalf("Log output should be valid JSON: %v", err)
	}

	// Check required fields
	if logEntry["level"] != "INFO" {
		t.Errorf("Expected level 'INFO', got '%v'", logEntry["level"])
	}

	if logEntry["type"] != "request" {
		t.Errorf("Expected type 'request', got '%v'", logEntry["type"])
	}

	if logEntry["method"] != "GET" {
		t.Errorf("Expected method 'GET', got '%v'", logEntry["method"])
	}

	if logEntry["path"] != "/api/users" {
		t.Errorf("Expected path '/api/users', got '%v'", logEntry["path"])
	}

	// Check timestamp exists
	if logEntry["timestamp"] == nil {
		t.Error("Expected timestamp field to be present")
	}

	// Check params
	paramsInterface, ok := logEntry["params"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected params to be a map")
	}

	if paramsInterface["id"] != "123" {
		t.Errorf("Expected param id '123', got '%v'", paramsInterface["id"])
	}
}

// TestLogResponse tests logging of HTTP response details
func TestLogResponse(t *testing.T) {
	logger := NewLogger()

	responseBody := map[string]interface{}{
		"message": "success",
		"data":    []string{"item1", "item2"},
	}

	output := captureOutput(func() {
		logger.LogResponse(200, responseBody)
	})

	// Verify output is valid JSON
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
	if err != nil {
		t.Fatalf("Log output should be valid JSON: %v", err)
	}

	// Check required fields
	if logEntry["level"] != "INFO" {
		t.Errorf("Expected level 'INFO', got '%v'", logEntry["level"])
	}

	if logEntry["type"] != "response" {
		t.Errorf("Expected type 'response', got '%v'", logEntry["type"])
	}

	statusCode, ok := logEntry["status_code"].(float64)
	if !ok || int(statusCode) != 200 {
		t.Errorf("Expected status_code 200, got '%v'", logEntry["status_code"])
	}

	// Check timestamp exists
	if logEntry["timestamp"] == nil {
		t.Error("Expected timestamp field to be present")
	}
}

// TestLogMatch tests logging when a rule is matched
func TestLogMatch(t *testing.T) {
	logger := NewLogger()

	rule := &models.MockRule{
		Path: "/api/test",
		Response: map[string]interface{}{
			"message": "test response",
		},
		Code: 201,
	}

	output := captureOutput(func() {
		logger.LogMatch(rule)
	})

	// Verify output is valid JSON
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
	if err != nil {
		t.Fatalf("Log output should be valid JSON: %v", err)
	}

	// Check required fields
	if logEntry["level"] != "INFO" {
		t.Errorf("Expected level 'INFO', got '%v'", logEntry["level"])
	}

	if logEntry["type"] != "match" {
		t.Errorf("Expected type 'match', got '%v'", logEntry["type"])
	}

	if logEntry["message"] != "Rule matched" {
		t.Errorf("Expected message 'Rule matched', got '%v'", logEntry["message"])
	}

	// Check rule details
	ruleInterface, ok := logEntry["rule"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected rule to be a map")
	}

	if ruleInterface["path"] != "/api/test" {
		t.Errorf("Expected rule path '/api/test', got '%v'", ruleInterface["path"])
	}

	code, ok := ruleInterface["code"].(float64)
	if !ok || int(code) != 201 {
		t.Errorf("Expected rule code 201, got '%v'", ruleInterface["code"])
	}
}

// TestLogDefault tests logging when default response is used
func TestLogDefault(t *testing.T) {
	logger := NewLogger()

	output := captureOutput(func() {
		logger.LogDefault()
	})

	// Verify output is valid JSON
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
	if err != nil {
		t.Fatalf("Log output should be valid JSON: %v", err)
	}

	// Check required fields
	if logEntry["level"] != "INFO" {
		t.Errorf("Expected level 'INFO', got '%v'", logEntry["level"])
	}

	if logEntry["type"] != "default" {
		t.Errorf("Expected type 'default', got '%v'", logEntry["type"])
	}

	expectedMessage := "No matching rule found, using default response"
	if logEntry["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%v'", expectedMessage, logEntry["message"])
	}

	// Check timestamp exists
	if logEntry["timestamp"] == nil {
		t.Error("Expected timestamp field to be present")
	}
}

// TestLogRequestWithEmptyParams tests logging request with empty parameters
func TestLogRequestWithEmptyParams(t *testing.T) {
	logger := NewLogger()

	output := captureOutput(func() {
		logger.LogRequest("POST", "/api/create", map[string]string{})
	})

	// Verify output is valid JSON
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
	if err != nil {
		t.Fatalf("Log output should be valid JSON: %v", err)
	}

	// Check that empty params are handled correctly
	paramsInterface, ok := logEntry["params"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected params to be a map")
	}

	if len(paramsInterface) != 0 {
		t.Errorf("Expected empty params map, got %v", paramsInterface)
	}
}

// TestLogResponseWithNilBody tests logging response with nil body
func TestLogResponseWithNilBody(t *testing.T) {
	logger := NewLogger()

	output := captureOutput(func() {
		logger.LogResponse(204, nil)
	})

	// Verify output is valid JSON
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
	if err != nil {
		t.Fatalf("Log output should be valid JSON: %v", err)
	}

	// Check that nil body is handled correctly
	if logEntry["body"] != nil {
		t.Errorf("Expected nil body, got '%v'", logEntry["body"])
	}
}

// TestLoggerWithUnmarshalableData tests handling of data that cannot be marshaled to JSON
func TestLoggerWithUnmarshalableData(t *testing.T) {
	logger := NewLogger()

	// Create a channel which cannot be marshaled to JSON
	unmarshalableData := make(chan int)

	output := captureOutput(func() {
		logger.LogResponse(200, unmarshalableData)
	})

	// Should output an error log instead of the original log
	if !strings.Contains(output, "log_error") {
		t.Error("Expected log_error when data cannot be marshaled")
	}

	if !strings.Contains(output, "Failed to marshal log entry") {
		t.Error("Expected error message about marshaling failure")
	}
}
