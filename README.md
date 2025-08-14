# Mock Service

A lightweight HTTP mock service built with Go and Gin framework. This service allows you to quickly set up mock HTTP endpoints by defining rules in a JSON configuration file.

## Features

- **Universal HTTP Handler**: Accepts requests for any path and HTTP method
- **JSON Configuration**: Define mock responses using simple JSON configuration files
- **Sequential Rule Matching**: Rules are processed in order, first match wins
- **Structured Logging**: All requests and responses are logged in JSON format
- **Query Parameter Support**: Automatically parses and logs query parameters
- **Docker Support**: Ready-to-use Docker configuration with multi-stage builds
- **Health Check Endpoint**: Built-in `/health` endpoint for monitoring
- **Graceful Shutdown**: Handles SIGINT and SIGTERM signals properly

## Quick Start

### Using Go

1. **Clone and build:**
   ```bash
   git clone <repository-url>
   cd mock-service
   go build -o bin/mock-service ./cmd/mock-service
   ```

2. **Create a configuration file:**
   ```json
   {
     "rules": [
       {
         "path": "/api/users",
         "response": {
           "users": ["alice", "bob"],
           "count": 2
         },
         "code": 200
       }
     ]
   }
   ```

3. **Run the service:**
   ```bash
   ./bin/mock-service -config config.json -port 8080
   ```

### Using Docker

1. **Build the Docker image:**
   ```bash
   docker build -t mock-service .
   ```

2. **Run with Docker:**
   ```bash
   docker run -p 8080:8080 -v $(pwd)/config:/app/config mock-service
   ```

### Using Docker Compose

1. **Start the service:**
   ```bash
   docker-compose up -d
   ```

## Configuration Format

The configuration file is a JSON document with the following structure:

```json
{
  "rules": [
    {
      "path": "/api/endpoint",
      "response": {
        "key": "value",
        "data": ["item1", "item2"]
      },
      "code": 200
    }
  ]
}
```

### Configuration Fields

- **`rules`** (array): List of mock rules to be processed
- **`path`** (string): The exact request path to match (case-sensitive)
- **`response`** (object): JSON response body to return when the rule matches
- **`code`** (integer): HTTP status code to return (defaults to 200 if not specified)

## Example Configurations

### Basic API Endpoints
```json
{
  "rules": [
    {
      "path": "/api/users",
      "response": {
        "users": [
          {"id": 1, "name": "Alice"},
          {"id": 2, "name": "Bob"}
        ]
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
}
```

### Error Responses
```json
{
  "rules": [
    {
      "path": "/api/error/500",
      "response": {
        "error": "Internal Server Error",
        "message": "Something went wrong"
      },
      "code": 500
    }
  ]
}
```

## Command Line Options

- **`-config`**: Path to configuration file (default: `config.json`)
- **`-port`**: Port to listen on (default: `8080`)

Example:
```bash
./mock-service -config /path/to/config.json -port 3000
```

## API Behavior

### Request Matching
- Rules are processed **sequentially** in the order they appear in the configuration
- **First matching rule wins** - subsequent rules are ignored
- Path matching is **case-sensitive** and requires **exact match**
- If no rule matches, returns a 404 "Not Found" response

### Query Parameters
- Query parameters are automatically parsed and logged
- Multiple values for the same parameter: only the first value is used
- Query parameters don't affect rule matching (only path is used for matching)

### HTTP Methods
- All HTTP methods are supported (GET, POST, PUT, DELETE, etc.)
- HTTP method doesn't affect rule matching (only path is used)

### Default Response
When no rule matches the request path:
```json
{}
```
Status Code: `200`

## Logging

All requests and responses are logged to stdout in JSON format:

### Request Log
```json
{
  "timestamp": "2024-01-14T15:30:45Z",
  "level": "INFO",
  "type": "request",
  "method": "GET",
  "path": "/api/users",
  "params": {"id": "123", "filter": "active"}
}
```

### Response Log
```json
{
  "timestamp": "2024-01-14T15:30:45Z",
  "level": "INFO",
  "type": "response",
  "status_code": 200,
  "body": {"users": ["alice", "bob"]}
}
```

### Rule Match Log
```json
{
  "timestamp": "2024-01-14T15:30:45Z",
  "level": "INFO",
  "type": "match",
  "message": "Rule matched",
  "rule": {"path": "/api/users", "code": 200}
}
```

## Docker Configuration

### Environment Variables
- **`CONFIG_FILE`**: Path to configuration file inside container (default: `/app/config/config.json`)

### Volume Mounts
Mount your configuration directory to `/app/config`:
```bash
docker run -v $(pwd)/config:/app/config mock-service
```

### Health Check
The service includes a health check endpoint at `/health`:
```json
{
  "status": "healthy",
  "service": "mock-service"
}
```

## Development

### Project Structure
```
mock-service/
├── cmd/mock-service/          # Main application entry point
├── internal/
│   ├── config/                # Configuration management
│   ├── handler/               # HTTP request handlers
│   ├── interfaces/            # Core interfaces
│   ├── logger/                # Logging functionality
│   ├── matcher/               # Path matching logic
│   ├── models/                # Data models
│   └── response/              # Response building
├── config/                    # Example configuration files
├── Dockerfile                 # Docker build configuration
├── docker-compose.yml         # Docker Compose configuration
└── README.md                  # This documentation
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...
```

### Building
```bash
# Build for current platform
go build -o bin/mock-service ./cmd/mock-service

# Build for Linux (for Docker)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mock-service ./cmd/mock-service
```

## Use Cases

- **API Development**: Mock external APIs during development
- **Testing**: Create predictable responses for automated tests
- **Prototyping**: Quickly prototype API responses without backend implementation
- **Integration Testing**: Simulate various API scenarios and error conditions
- **Load Testing**: Test client applications against consistent mock responses

## License

This project is open source. Please check the license file for details.