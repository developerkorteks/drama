#!/bin/bash

# DramaQu API Production Deployment Script
# Usage: ./deploy.sh [command]
# Commands: build, deploy, restart, stop, logs, status, clean

set -e

# Configuration
APP_NAME="dramaqu-api"
IMAGE_NAME="dramaqu-api:latest"
CONTAINER_NAME="dramaqu-api-prod"
PORT="52983"
COMPOSE_FILE="docker-compose.prod.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        error "Docker is not running. Please start Docker first."
    fi
}

# Check if Docker Compose is available
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed. Please install Docker Compose first."
    fi
}

# Generate Swagger documentation
generate_swagger() {
    log "Generating Swagger documentation..."
    if command -v swag &> /dev/null; then
        swag init --parseDependency --parseInternal
        log "Swagger documentation generated successfully"
    else
        warn "Swag CLI not found. Swagger documentation may be outdated."
    fi
}

# Build Docker image
build_image() {
    log "Building Docker image: $IMAGE_NAME"
    
    # Generate swagger docs first
    generate_swagger
    
    # Build the image
    docker build -t $IMAGE_NAME . || error "Failed to build Docker image"
    
    log "Docker image built successfully: $IMAGE_NAME"
}

# Deploy application
deploy_app() {
    log "Deploying $APP_NAME to production..."
    
    # Check prerequisites
    check_docker
    check_docker_compose
    
    # Build image
    build_image
    
    # Stop existing container if running
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        log "Stopping existing container..."
        docker-compose -f $COMPOSE_FILE down
    fi
    
    # Start new container
    log "Starting new container..."
    docker-compose -f $COMPOSE_FILE up -d
    
    # Wait for container to be ready
    log "Waiting for application to start..."
    sleep 10
    
    # Health check
    if health_check; then
        log "✅ Deployment successful!"
        log "🌐 Application is running on port $PORT"
        log "📚 Swagger documentation: http://[your-domain]/swagger/index.html"
        log "🔧 Dynamic host detection enabled - supports wildcard domains"
    else
        error "❌ Deployment failed - health check failed"
    fi
}

# Restart application
restart_app() {
    log "Restarting $APP_NAME..."
    check_docker_compose
    
    docker-compose -f $COMPOSE_FILE restart
    
    sleep 5
    if health_check; then
        log "✅ Application restarted successfully"
    else
        error "❌ Restart failed - health check failed"
    fi
}

# Stop application
stop_app() {
    log "Stopping $APP_NAME..."
    check_docker_compose
    
    docker-compose -f $COMPOSE_FILE down
    log "✅ Application stopped"
}

# Show logs
show_logs() {
    log "Showing logs for $APP_NAME..."
    check_docker_compose
    
    docker-compose -f $COMPOSE_FILE logs -f --tail=100
}

# Show status
show_status() {
    log "Checking status of $APP_NAME..."
    
    echo -e "\n${BLUE}=== Docker Images ===${NC}"
    docker images | grep -E "(REPOSITORY|$APP_NAME)" || echo "No images found"
    
    echo -e "\n${BLUE}=== Running Containers ===${NC}"
    docker ps | grep -E "(CONTAINER|$CONTAINER_NAME)" || echo "No containers running"
    
    echo -e "\n${BLUE}=== Docker Compose Status ===${NC}"
    if [ -f "$COMPOSE_FILE" ]; then
        docker-compose -f $COMPOSE_FILE ps
    else
        echo "Docker compose file not found: $COMPOSE_FILE"
    fi
    
    echo -e "\n${BLUE}=== Health Check ===${NC}"
    if health_check; then
        echo -e "${GREEN}✅ Application is healthy${NC}"
    else
        echo -e "${RED}❌ Application is not responding${NC}"
    fi
}

# Health check
health_check() {
    local max_attempts=5
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s http://localhost:$PORT/swagger-config > /dev/null 2>&1; then
            return 0
        fi
        
        log "Health check attempt $attempt/$max_attempts failed, retrying..."
        sleep 2
        ((attempt++))
    done
    
    return 1
}

# Clean up
clean_up() {
    log "Cleaning up Docker resources..."
    
    # Stop and remove containers
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        docker-compose -f $COMPOSE_FILE down -v
    fi
    
    # Remove image
    if docker images -q $IMAGE_NAME | grep -q .; then
        docker rmi $IMAGE_NAME 2>/dev/null || warn "Could not remove image $IMAGE_NAME"
    fi
    
    # Clean up unused resources
    docker system prune -f
    
    log "✅ Cleanup completed"
}

# Backup function
backup_data() {
    log "Creating backup..."
    
    local backup_dir="backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p $backup_dir
    
    # Backup configuration files
    cp -r .env* $backup_dir/ 2>/dev/null || true
    cp docker-compose*.yml $backup_dir/ 2>/dev/null || true
    cp Dockerfile $backup_dir/ 2>/dev/null || true
    
    log "✅ Backup created in $backup_dir"
}

# Update application
update_app() {
    log "Updating $APP_NAME..."
    
    # Create backup first
    backup_data
    
    # Pull latest code (if using git)
    if [ -d ".git" ]; then
        log "Pulling latest code from git..."
        git pull origin main || warn "Git pull failed or not in git repository"
    fi
    
    # Rebuild and deploy
    deploy_app
}

# Show help
show_help() {
    echo -e "${BLUE}DramaQu API Production Deployment Script${NC}"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  build     - Build Docker image"
    echo "  deploy    - Deploy application to production"
    echo "  restart   - Restart the application"
    echo "  stop      - Stop the application"
    echo "  logs      - Show application logs"
    echo "  status    - Show application status"
    echo "  clean     - Clean up Docker resources"
    echo "  backup    - Create backup of configuration"
    echo "  update    - Update and redeploy application"
    echo "  health    - Check application health"
    echo "  help      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 deploy    # Deploy to production"
    echo "  $0 logs      # Show logs"
    echo "  $0 status    # Check status"
    echo ""
}

# Main script logic
case "${1:-help}" in
    "build")
        build_image
        ;;
    "deploy")
        deploy_app
        ;;
    "restart")
        restart_app
        ;;
    "stop")
        stop_app
        ;;
    "logs")
        show_logs
        ;;
    "status")
        show_status
        ;;
    "clean")
        clean_up
        ;;
    "backup")
        backup_data
        ;;
    "update")
        update_app
        ;;
    "health")
        if health_check; then
            log "✅ Application is healthy"
        else
            error "❌ Application is not healthy"
        fi
        ;;
    "help"|*)
        show_help
        ;;
esac