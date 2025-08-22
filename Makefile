# Variables
APP_NAME=dramaqu-api
DOCKER_IMAGE=dramaqu-api:latest
PORT=52983

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build       - Build the Docker image"
	@echo "  run         - Run the application with Docker Compose"
	@echo "  stop        - Stop the application"
	@echo "  restart     - Restart the application"
	@echo "  logs        - Show application logs"
	@echo "  clean       - Clean up Docker resources"
	@echo "  dev         - Run in development mode"
	@echo "  prod        - Run in production mode"
	@echo "  swagger     - Generate Swagger documentation"
	@echo "  test        - Run tests"

# Build Docker image
.PHONY: build
build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run with Docker Compose
.PHONY: run
run:
	@echo "Starting application..."
	docker-compose up -d

# Stop application
.PHONY: stop
stop:
	@echo "Stopping application..."
	docker-compose down

# Restart application
.PHONY: restart
restart: stop run

# Show logs
.PHONY: logs
logs:
	docker-compose logs -f

# Clean up Docker resources
.PHONY: clean
clean:
	@echo "Cleaning up Docker resources..."
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	docker system prune -f

# Development mode
.PHONY: dev
dev:
	@echo "Running in development mode..."
	GIN_MODE=debug PORT=$(PORT) go run main.go

# Production mode
.PHONY: prod
prod:
	@echo "Running in production mode..."
	GIN_MODE=release PORT=$(PORT) go run main.go

# Generate Swagger documentation
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	swag init --parseDependency --parseInternal

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Build and run
.PHONY: deploy
deploy: build run

# Health check
.PHONY: health
health:
	@echo "Checking application health..."
	curl -f http://localhost:$(PORT)/health || echo "Application is not healthy"

# Show application status
.PHONY: status
status:
	@echo "Application status:"
	docker-compose ps