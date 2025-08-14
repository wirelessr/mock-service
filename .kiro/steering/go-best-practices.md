---
inclusion: always
---

# Go Development Best Practices

## Code Style and Structure

### Package Organization
- Use clear, descriptive package names (avoid generic names like `util`, `common`)
- Keep packages focused on a single responsibility
- Place interfaces close to their usage, not their implementation
- Use internal packages for code that shouldn't be imported externally

### Naming Conventions
- Use camelCase for variables and functions
- Use PascalCase for exported types and functions
- Use descriptive names that explain purpose, not implementation
- Prefer longer, clear names over short, cryptic ones
- Use consistent naming patterns across the codebase

### Error Handling
- Always handle errors explicitly - never ignore them
- Use `fmt.Errorf` with `%w` verb for error wrapping
- Create custom error types for domain-specific errors
- Return errors as the last return value
- Use early returns to reduce nesting

### Testing
- Write table-driven tests for multiple test cases
- Use descriptive test function names that explain what is being tested
- Test both success and failure scenarios
- Aim for high test coverage (80%+) but focus on meaningful tests
- Use testify/assert for cleaner test assertions when appropriate
- Mock external dependencies using interfaces

## HTTP Service Development

### Gin Framework Usage
- Use middleware for cross-cutting concerns (logging, authentication, CORS)
- Group related routes using router groups
- Use proper HTTP status codes (200, 201, 400, 404, 500, etc.)
- Validate input data before processing
- Use JSON binding for request parsing

### API Design
- Follow RESTful conventions where applicable
- Use consistent response formats
- Include proper error messages in API responses
- Implement health check endpoints
- Use structured logging for better observability

### Configuration Management
- Use environment variables for configuration
- Provide sensible defaults
- Validate configuration on startup
- Support configuration files (JSON, YAML) when appropriate

## Docker Best Practices

### Dockerfile Optimization
- Use multi-stage builds to reduce image size
- Use Alpine Linux as base image for smaller footprint
- Run as non-root user for security
- Copy only necessary files
- Use .dockerignore to exclude unnecessary files

### Container Security
- Create dedicated user accounts (avoid root)
- Use specific version tags, not `latest`
- Minimize attack surface by including only required dependencies
- Implement health checks

## Code Quality

### Performance Considerations
- Use buffered channels when appropriate
- Avoid premature optimization
- Profile code when performance issues arise
- Use sync.Pool for frequently allocated objects
- Be mindful of memory allocations in hot paths

### Concurrency
- Use channels for communication between goroutines
- Prefer channels over shared memory with mutexes
- Use context.Context for cancellation and timeouts
- Be careful with goroutine lifecycle management
- Use sync package primitives when channels aren't suitable

### Documentation
- Write clear, concise comments for exported functions and types
- Use godoc format for documentation comments
- Include examples in documentation when helpful
- Keep README files up to date
- Document complex algorithms and business logic

## Development Workflow

### Git Practices
- Use conventional commit messages (feat:, fix:, docs:, etc.)
- Keep commits atomic and focused
- Write descriptive commit messages
- Use feature branches for development
- Squash commits when merging to keep history clean

### Code Review
- Review for correctness, readability, and maintainability
- Check for proper error handling
- Verify test coverage for new code
- Ensure documentation is updated
- Look for potential security issues