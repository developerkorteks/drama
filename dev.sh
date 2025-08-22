#!/bin/bash

# Development Script for DramaQu API
# Usage: ./dev.sh [command]

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
PORT="52983"
DEV_COMPOSE_FILE="docker-compose.yml"

log() {
    echo -e "${GREEN}[DEV] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[DEV] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[DEV] ERROR: $1${NC}"
    exit 1
}

# Generate Swagger docs
generate_docs() {
    log "Generating Swagger documentation..."
    if command -v swag &> /dev/null; then
        swag init --parseDependency --parseInternal
        log "✅ Swagger docs generated"
    else
        warn "Swag CLI not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"
    fi
}

# Run in development mode
run_dev() {
    log "Starting development server..."
    
    # Generate docs first
    generate_docs
    
    # Set development environment
    export GIN_MODE=debug
    export PORT=$PORT
    
    log "🚀 Starting server on http://localhost:$PORT"
    log "📚 Swagger docs: http://localhost:$PORT/swagger/index.html"
    
    # Run the application
    go run main.go
}

# Run with Docker (development)
run_docker_dev() {
    log "Starting development server with Docker..."
    
    generate_docs
    
    # Use development compose file
    docker-compose -f $DEV_COMPOSE_FILE up --build
}

# Run tests
run_tests() {
    log "Running tests..."
    
    if [ -d "tests" ] || find . -name "*_test.go" -type f | grep -q .; then
        go test -v ./...
    else
        warn "No tests found"
    fi
}

# Install dependencies
install_deps() {
    log "Installing dependencies..."
    
    # Go dependencies
    go mod tidy
    go mod download
    
    # Install swag if not present
    if ! command -v swag &> /dev/null; then
        log "Installing Swag CLI..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    
    log "✅ Dependencies installed"
}

# Format code
format_code() {
    log "Formatting code..."
    
    go fmt ./...
    
    if command -v goimports &> /dev/null; then
        goimports -w .
    else
        warn "goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"
    fi
    
    log "✅ Code formatted"
}

# Lint code
lint_code() {
    log "Linting code..."
    
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
    else
        warn "golangci-lint not found. Install from: https://golangci-lint.run/usage/install/"
        
        # Fallback to basic checks
        go vet ./...
    fi
}

# Clean development environment
clean_dev() {
    log "Cleaning development environment..."
    
    # Stop docker containers
    docker-compose -f $DEV_COMPOSE_FILE down 2>/dev/null || true
    
    # Clean Go cache
    go clean -cache
    go clean -modcache
    
    # Remove generated docs (optional)
    # rm -rf docs/
    
    log "✅ Development environment cleaned"
}

# Show development info
show_info() {
    echo -e "${BLUE}=== DramaQu API Development Environment ===${NC}"
    echo ""
    echo "Go version: $(go version)"
    echo "Port: $PORT"
    echo "Mode: Development"
    echo ""
    echo -e "${BLUE}=== Available Endpoints ===${NC}"
    echo "• Health: http://localhost:$PORT/swagger-config"
    echo "• Swagger: http://localhost:$PORT/swagger/index.html"
    echo "• Docs: http://localhost:$PORT/docs/index.html"
    echo ""
    echo -e "${BLUE}=== Project Structure ===${NC}"
    find . -type f -name "*.go" | head -10
    echo ""
}

# Hot reload (if air is installed)
hot_reload() {
    if command -v air &> /dev/null; then
        log "Starting hot reload with Air..."
        air
    else
        warn "Air not found. Install with: go install github.com/cosmtrek/air@latest"
        log "Falling back to normal development mode..."
        run_dev
    fi
}

# Show help
show_help() {
    echo -e "${BLUE}DramaQu API Development Script${NC}"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  run       - Run development server"
    echo "  docker    - Run with Docker (development)"
    echo "  test      - Run tests"
    echo "  install   - Install dependencies"
    echo "  docs      - Generate Swagger documentation"
    echo "  format    - Format code"
    echo "  lint      - Lint code"
    echo "  clean     - Clean development environment"
    echo "  info      - Show development info"
    echo "  hot       - Run with hot reload (requires air)"
    echo "  help      - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 run      # Start development server"
    echo "  $0 docs     # Generate Swagger docs"
    echo "  $0 test     # Run tests"
    echo ""
}

# Main logic
case "${1:-help}" in
    "run")
        run_dev
        ;;
    "docker")
        run_docker_dev
        ;;
    "test")
        run_tests
        ;;
    "install")
        install_deps
        ;;
    "docs")
        generate_docs
        ;;
    "format")
        format_code
        ;;
    "lint")
        lint_code
        ;;
    "clean")
        clean_dev
        ;;
    "info")
        show_info
        ;;
    "hot")
        hot_reload
        ;;
    "help"|*)
        show_help
        ;;
esac