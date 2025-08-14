package matcher

import "mock-service/internal/models"

// PathMatcherImpl implements the PathMatcher interface
// It handles matching request paths against configured rules
type PathMatcherImpl struct{}

// NewPathMatcher creates a new instance of PathMatcher
func NewPathMatcher() *PathMatcherImpl {
	return &PathMatcherImpl{}
}

// FindMatch finds the first matching rule for the given request path
// Returns the matched rule and true if found, nil and false otherwise
// Rules are processed sequentially in the order they appear in the configuration
func (pm *PathMatcherImpl) FindMatch(requestPath string, rules []models.MockRule) (*models.MockRule, bool) {
	// Handle empty rules list
	if len(rules) == 0 {
		return nil, false
	}

	// Iterate through rules sequentially to find exact match
	for i := range rules {
		if rules[i].Path == requestPath {
			return &rules[i], true
		}
	}

	// No match found
	return nil, false
}
