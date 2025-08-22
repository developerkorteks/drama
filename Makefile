.PHONY: run build docs clean test

# Run the application
run:
	go run main.go

# Build the application
build:
	go build -o bin/dramaqu main.go

# Generate Swagger documentation
docs:
	swag init

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Test the application
test:
	go test ./...

# Run with hot reload (requires air)
dev:
	air

# Install air for hot reload
install-air:
	go install github.com/cosmtrek/air@latest

# Setup project (install deps and generate docs)
setup: deps docs
	@echo "Project setup complete!"

# Run in production mode
prod:
	GIN_MODE=release go run main.go