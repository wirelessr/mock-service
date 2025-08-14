---
inclusion: always
---

# Mock Service Development Context

## Project Overview

This is a lightweight HTTP mock service built with Go and Gin framework. The service accepts any HTTP request path and method, returning configured responses based on JSON configuration files.

## Architecture Principles

### Clean Architecture
- Separate business logic from HTTP handling
- Use interfaces to define contracts between layers
- Keep external dependencies (Gin, JSON parsing) isolated
- Make components easily testable through dependency injection

### Component Structure
- **Models**: Data structures (MockRule, Config)
- **Interfaces**: Contracts for all major components
- **Config**: Configuration file loading and management
- **Matcher**: Path matching logic
- **Response**: HTTP response building
- **Logger**: Structured logging functionality
- **Handler**: HTTP request handling (Gin integration)

## Key Requirements

### Functional Requirements
- Accept any HTTP path and method (universal handler)
- Load configuration from JSON files
- Match requests against configured rules sequentially
- Return 200 status with empty JSON `{}` for unmatched paths
- Support custom status codes and response bodies
- Log all requests and responses in structured JSON format

### Non-Functional Requirements
- Lightweight and fast startup
- Docker-ready with small image size (<50MB)
- High test coverage (80%+)
- Clear error messages and logging
- Graceful shutdown handling

## Configuration Format

```json
{
  "rules": [
    {
      "path": "/api/endpoint",
      "response": { "key": "value" },
      "code": 200
    }
  ]
}
```

## Development Guidelines

### Testing Strategy
- Unit tests for all business logic components
- Integration tests for HTTP handlers
- Mock external dependencies using interfaces
- Test both success and error scenarios
- Maintain high test coverage

### Error Handling
- Configuration loading errors should prevent startup
- Invalid JSON should return clear error messages
- Runtime errors should be logged but not crash the service
- Use structured logging for better observability

### Performance Considerations
- Sequential rule matching (first match wins)
- Efficient JSON parsing and response building
- Minimal memory allocations in request handling
- Fast startup time for container environments

## Docker Deployment

### Container Requirements
- Multi-stage build for minimal image size
- Alpine Linux base image
- Non-root user execution
- Health check endpoint at `/health`
- Configuration file mounting at `/app/config/`

### Environment Variables
- `CONFIG_FILE`: Path to configuration file
- Port configuration via command line flags

## Logging Format

All logs use structured JSON format:
- Request logs: method, path, query parameters
- Response logs: status code, response body
- Match logs: which rule was matched
- Default logs: when no rule matches

## Common Patterns

### Component Initialization
```go
// Create all components
configManager := config.NewConfigManager()
pathMatcher := matcher.NewPathMatcher()
responseBuilder := response.NewResponseBuilder()
logger := logger.NewLogger()

// Wire them together
handler := handler.NewUniversalHandler(
    configManager, pathMatcher, responseBuilder, logger
)
```

### Error Handling Pattern
```go
if err := operation(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### Testing Pattern
```go
func TestComponent(t *testing.T) {
    // Arrange
    component := NewComponent()
    
    // Act
    result, err := component.DoSomething()
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // ... assertions
}
```