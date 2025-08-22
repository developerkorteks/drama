#!/bin/bash

# Monitoring Script for DramaQu API
# Usage: ./monitor.sh [command]

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
PORT="52983"
CONTAINER_NAME="dramaqu-api-prod"
LOG_FILE="/tmp/dramaqu-monitor.log"

log() {
    echo -e "${GREEN}[MONITOR] $1${NC}"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >> $LOG_FILE
}

warn() {
    echo -e "${YELLOW}[MONITOR] WARNING: $1${NC}"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1" >> $LOG_FILE
}

error() {
    echo -e "${RED}[MONITOR] ERROR: $1${NC}"
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1" >> $LOG_FILE
}

# Health check function
health_check() {
    local endpoint="http://localhost:$PORT/swagger-config"
    
    if curl -f -s --max-time 10 "$endpoint" > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Check container status
check_container() {
    if docker ps --format "table {{.Names}}\t{{.Status}}" | grep -q "$CONTAINER_NAME"; then
        return 0
    else
        return 1
    fi
}

# Get container stats
get_container_stats() {
    if check_container; then
        docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}" $CONTAINER_NAME
    else
        echo "Container not running"
    fi
}

# Monitor continuously
monitor_continuous() {
    log "Starting continuous monitoring..."
    log "Press Ctrl+C to stop monitoring"
    
    local check_interval=30
    local failure_count=0
    local max_failures=3
    
    while true; do
        echo -e "\n${BLUE}=== Health Check $(date +'%Y-%m-%d %H:%M:%S') ===${NC}"
        
        if health_check; then
            echo -e "${GREEN}✅ API is healthy${NC}"
            failure_count=0
        else
            failure_count=$((failure_count + 1))
            echo -e "${RED}❌ API health check failed (${failure_count}/${max_failures})${NC}"
            
            if [ $failure_count -ge $max_failures ]; then
                error "API has failed $max_failures consecutive health checks!"
                
                # Try to restart if container is running
                if check_container; then
                    warn "Attempting to restart container..."
                    docker restart $CONTAINER_NAME
                    sleep 10
                    failure_count=0
                else
                    error "Container is not running!"
                fi
            fi
        fi
        
        # Show container status
        echo -e "\n${BLUE}=== Container Status ===${NC}"
        if check_container; then
            echo -e "${GREEN}✅ Container is running${NC}"
            get_container_stats
        else
            echo -e "${RED}❌ Container is not running${NC}"
        fi
        
        # Show recent logs
        echo -e "\n${BLUE}=== Recent Logs ===${NC}"
        if check_container; then
            docker logs --tail=5 $CONTAINER_NAME 2>/dev/null || echo "No logs available"
        fi
        
        sleep $check_interval
    done
}

# Quick status check
quick_status() {
    echo -e "${BLUE}=== DramaQu API Status ===${NC}"
    
    # API Health
    echo -n "API Health: "
    if health_check; then
        echo -e "${GREEN}✅ Healthy${NC}"
    else
        echo -e "${RED}❌ Unhealthy${NC}"
    fi
    
    # Container Status
    echo -n "Container: "
    if check_container; then
        echo -e "${GREEN}✅ Running${NC}"
    else
        echo -e "${RED}❌ Not Running${NC}"
    fi
    
    # Port Status
    echo -n "Port $PORT: "
    if netstat -tuln 2>/dev/null | grep -q ":$PORT "; then
        echo -e "${GREEN}✅ Open${NC}"
    else
        echo -e "${RED}❌ Closed${NC}"
    fi
    
    # Disk Space
    echo -n "Disk Space: "
    local disk_usage=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ "$disk_usage" -lt 80 ]; then
        echo -e "${GREEN}✅ ${disk_usage}% used${NC}"
    elif [ "$disk_usage" -lt 90 ]; then
        echo -e "${YELLOW}⚠️  ${disk_usage}% used${NC}"
    else
        echo -e "${RED}❌ ${disk_usage}% used${NC}"
    fi
    
    # Memory Usage
    echo -n "Memory: "
    local mem_usage=$(free | awk 'NR==2{printf "%.1f", $3*100/$2}')
    echo -e "${GREEN}${mem_usage}% used${NC}"
}

# Show detailed stats
detailed_stats() {
    echo -e "${BLUE}=== Detailed Statistics ===${NC}"
    
    # System info
    echo -e "\n${BLUE}System Information:${NC}"
    echo "Hostname: $(hostname)"
    echo "Uptime: $(uptime -p)"
    echo "Load Average: $(uptime | awk -F'load average:' '{print $2}')"
    
    # Docker info
    echo -e "\n${BLUE}Docker Information:${NC}"
    if check_container; then
        echo "Container Status: Running"
        echo "Container Uptime: $(docker inspect --format='{{.State.StartedAt}}' $CONTAINER_NAME)"
        
        echo -e "\n${BLUE}Resource Usage:${NC}"
        get_container_stats
        
        echo -e "\n${BLUE}Container Logs (last 20 lines):${NC}"
        docker logs --tail=20 $CONTAINER_NAME
    else
        echo "Container Status: Not Running"
    fi
    
    # Network info
    echo -e "\n${BLUE}Network Information:${NC}"
    echo "Listening on port $PORT: $(netstat -tuln 2>/dev/null | grep ":$PORT " || echo "No")"
    
    # API endpoints test
    echo -e "\n${BLUE}API Endpoints Test:${NC}"
    local endpoints=("/swagger-config" "/swagger/index.html")
    
    for endpoint in "${endpoints[@]}"; do
        echo -n "Testing $endpoint: "
        if curl -f -s --max-time 5 "http://localhost:$PORT$endpoint" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ OK${NC}"
        else
            echo -e "${RED}❌ Failed${NC}"
        fi
    done
}

# Show logs with filtering
show_logs() {
    local lines=${2:-100}
    local filter=${3:-""}
    
    echo -e "${BLUE}=== Application Logs (last $lines lines) ===${NC}"
    
    if check_container; then
        if [ -n "$filter" ]; then
            docker logs --tail=$lines $CONTAINER_NAME 2>&1 | grep -i "$filter"
        else
            docker logs --tail=$lines $CONTAINER_NAME
        fi
    else
        echo "Container is not running"
    fi
}

# Performance test
performance_test() {
    echo -e "${BLUE}=== Performance Test ===${NC}"
    
    if ! command -v curl &> /dev/null; then
        error "curl is required for performance testing"
        return 1
    fi
    
    local endpoint="http://localhost:$PORT/swagger-config"
    local requests=10
    
    echo "Testing $endpoint with $requests requests..."
    
    local total_time=0
    local successful_requests=0
    
    for i in $(seq 1 $requests); do
        local start_time=$(date +%s.%N)
        
        if curl -f -s --max-time 10 "$endpoint" > /dev/null 2>&1; then
            local end_time=$(date +%s.%N)
            local request_time=$(echo "$end_time - $start_time" | bc -l)
            total_time=$(echo "$total_time + $request_time" | bc -l)
            successful_requests=$((successful_requests + 1))
            echo "Request $i: ${request_time}s"
        else
            echo "Request $i: Failed"
        fi
    done
    
    if [ $successful_requests -gt 0 ]; then
        local avg_time=$(echo "scale=3; $total_time / $successful_requests" | bc -l)
        echo -e "\n${GREEN}Results:${NC}"
        echo "Successful requests: $successful_requests/$requests"
        echo "Average response time: ${avg_time}s"
        echo "Success rate: $(echo "scale=1; $successful_requests * 100 / $requests" | bc -l)%"
    else
        echo -e "\n${RED}All requests failed${NC}"
    fi
}

# Show help
show_help() {
    echo -e "${BLUE}DramaQu API Monitoring Script${NC}"
    echo ""
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "Commands:"
    echo "  status    - Quick status check"
    echo "  monitor   - Continuous monitoring (Ctrl+C to stop)"
    echo "  stats     - Detailed statistics"
    echo "  logs      - Show application logs"
    echo "  perf      - Run performance test"
    echo "  help      - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 status           # Quick status check"
    echo "  $0 monitor          # Start continuous monitoring"
    echo "  $0 logs 50          # Show last 50 log lines"
    echo "  $0 logs 100 error   # Show last 100 lines containing 'error'"
    echo ""
}

# Main logic
case "${1:-help}" in
    "status")
        quick_status
        ;;
    "monitor")
        monitor_continuous
        ;;
    "stats")
        detailed_stats
        ;;
    "logs")
        show_logs "$@"
        ;;
    "perf")
        performance_test
        ;;
    "help"|*)
        show_help
        ;;
esac