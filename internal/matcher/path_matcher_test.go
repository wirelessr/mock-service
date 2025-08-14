package matcher

import (
	"testing"

	"mock-service/internal/models"
)

// TestNewPathMatcher tests the creation of a new PathMatcher instance
func TestNewPathMatcher(t *testing.T) {
	pm := NewPathMatcher()
	if pm == nil {
		t.Fatal("NewPathMatcher should return a non-nil instance")
	}
}

// TestFindMatchExactPath tests exact path matching functionality
func TestFindMatchExactPath(t *testing.T) {
	pm := NewPathMatcher()

	rules := []models.MockRule{
		{
			Path:     "/api/users",
			Response: map[string]interface{}{"message": "Users endpoint"},
			Code:     200,
		},
		{
			Path:     "/api/products",
			Response: map[string]interface{}{"message": "Products endpoint"},
			Code:     200,
		},
	}

	// Test exact match for first rule
	rule, found := pm.FindMatch("/api/users", rules)
	if !found {
		t.Error("Expected to find match for '/api/users'")
	}
	if rule == nil {
		t.Fatal("Expected non-nil rule")
	}
	if rule.Path != "/api/users" {
		t.Errorf("Expected path '/api/users', got '%s'", rule.Path)
	}

	// Test exact match for second rule
	rule, found = pm.FindMatch("/api/products", rules)
	if !found {
		t.Error("Expected to find match for '/api/products'")
	}
	if rule == nil {
		t.Fatal("Expected non-nil rule")
	}
	if rule.Path != "/api/products" {
		t.Errorf("Expected path '/api/products', got '%s'", rule.Path)
	}
}

// TestFindMatchNoMatch tests behavior when no rule matches the request path
func TestFindMatchNoMatch(t *testing.T) {
	pm := NewPathMatcher()

	rules := []models.MockRule{
		{
			Path:     "/api/users",
			Response: map[string]interface{}{"message": "Users endpoint"},
			Code:     200,
		},
	}

	// Test no match scenario
	rule, found := pm.FindMatch("/api/orders", rules)
	if found {
		t.Error("Expected no match for '/api/orders'")
	}
	if rule != nil {
		t.Error("Expected nil rule when no match found")
	}
}

// TestFindMatchEmptyRules tests behavior with empty rules list
func TestFindMatchEmptyRules(t *testing.T) {
	pm := NewPathMatcher()

	var rules []models.MockRule

	// Test with empty rules
	rule, found := pm.FindMatch("/any/path", rules)
	if found {
		t.Error("Expected no match with empty rules")
	}
	if rule != nil {
		t.Error("Expected nil rule with empty rules")
	}
}

// TestFindMatchSequentialMatching tests that rules are processed in order
func TestFindMatchSequentialMatching(t *testing.T) {
	pm := NewPathMatcher()

	rules := []models.MockRule{
		{
			Path:     "/api/test",
			Response: map[string]interface{}{"message": "First rule"},
			Code:     200,
		},
		{
			Path:     "/api/test",
			Response: map[string]interface{}{"message": "Second rule"},
			Code:     404,
		},
	}

	// Should match the first rule (sequential processing)
	rule, found := pm.FindMatch("/api/test", rules)
	if !found {
		t.Error("Expected to find match for '/api/test'")
	}
	if rule == nil {
		t.Fatal("Expected non-nil rule")
	}

	// Verify it matched the first rule, not the second
	if rule.Code != 200 {
		t.Errorf("Expected first rule (code 200), got code %d", rule.Code)
	}

	response, ok := rule.Response["message"].(string)
	if !ok {
		t.Fatal("Expected string message in response")
	}
	if response != "First rule" {
		t.Errorf("Expected 'First rule', got '%s'", response)
	}
}

// TestFindMatchCaseSensitive tests that path matching is case sensitive
func TestFindMatchCaseSensitive(t *testing.T) {
	pm := NewPathMatcher()

	rules := []models.MockRule{
		{
			Path:     "/api/Users",
			Response: map[string]interface{}{"message": "Users endpoint"},
			Code:     200,
		},
	}

	// Test case sensitivity - should not match
	rule, found := pm.FindMatch("/api/users", rules)
	if found {
		t.Error("Expected no match for case-different path '/api/users'")
	}
	if rule != nil {
		t.Error("Expected nil rule for case-different path")
	}

	// Test exact case - should match
	rule, found = pm.FindMatch("/api/Users", rules)
	if !found {
		t.Error("Expected match for exact case '/api/Users'")
	}
	if rule == nil {
		t.Error("Expected non-nil rule for exact case match")
	}
}
