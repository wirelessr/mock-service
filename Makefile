.PHONY: help build test clean lint fmt vet coverage docker-build docker-test run

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the binary
	@echo "Building mock-service..."
	@CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/mock-service ./cmd/mock-service

build-all: ## Build binaries for all platforms
	@echo "Building for Linux..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/mock-service-linux ./cmd/mock-service
	@echo "Building for macOS..."
	@CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o bin/mock-service-darwin ./cmd/mock-service
	@echo "Building for Windows..."
	@CGO_ENABLED=0 GOOS=windows go build -a -installsuffix cgo -o bin/mock-service.exe ./cmd/mock-service

# Test targets
test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out

# Code quality targets
fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@golangci-lint run

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t mock-service:latest .

docker-test: docker-build ## Test Docker image
	@echo "Testing Docker image..."
	@mkdir -p test-config
	@echo '{"rules":[{"path":"/test","response":{"message":"docker test"},"code":200}]}' > test-config/config.json
	@docker run -d --name mock-service-test -p 8080:8080 -v $(PWD)/test-config:/app/config mock-service:latest ./mock-service -config /app/config/config.json -port 8080
	@sleep 3
	@curl -f http://localhost:8080/health || (docker logs mock-service-test && docker stop mock-service-test && docker rm mock-service-test && exit 1)
	@curl -f http://localhost:8080/test || (docker logs mock-service-test && docker stop mock-service-test && docker rm mock-service-test && exit 1)
	@docker stop mock-service-test
	@docker rm mock-service-test
	@rm -rf test-config
	@echo "Docker test passed!"

# Development targets
run: build ## Build and run the service
	@echo "Starting mock-service..."
	@./bin/mock-service -config config/example-basic.json -port 8080

dev: ## Run in development mode with example config
	@echo "Starting mock-service in development mode..."
	@go run ./cmd/mock-service -config config/example-basic.json -port 8080

# Cleanup targets
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -rf test-config/

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# CI targets (used by GitHub Actions)
ci-test: fmt vet lint test-coverage ## Run all CI checks

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest