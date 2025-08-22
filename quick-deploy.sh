#!/bin/bash

# Quick Deploy Script for DramaQu API
# One-command deployment for production

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔══════════════════════════════════════╗"
echo "║        DramaQu API Quick Deploy      ║"
echo "║              Production              ║"
echo "╚══════════════════════════════════════╝"
echo -e "${NC}"

# Check if running as root (for production servers)
if [[ $EUID -eq 0 ]]; then
   echo -e "${YELLOW}Warning: Running as root. Consider using a non-root user for security.${NC}"
fi

# Make deploy script executable
chmod +x deploy.sh

echo -e "${GREEN}🚀 Starting quick deployment...${NC}"

# Run full deployment
./deploy.sh deploy

echo -e "${GREEN}"
echo "╔══════════════════════════════════════╗"
echo "║         Deployment Complete!        ║"
echo "╚══════════════════════════════════════╝"
echo -e "${NC}"

echo -e "${BLUE}📋 Next Steps:${NC}"
echo "1. Configure your reverse proxy (Nginx/Traefik) to point to port 52983"
echo "2. Set up SSL certificate for HTTPS"
echo "3. Configure your wildcard domain DNS"
echo "4. Test the API endpoints"
echo ""
echo -e "${BLUE}🔗 Useful Commands:${NC}"
echo "  ./deploy.sh status  - Check application status"
echo "  ./deploy.sh logs    - View application logs"
echo "  ./deploy.sh restart - Restart application"
echo "  ./deploy.sh stop    - Stop application"
echo ""
echo -e "${GREEN}✅ DramaQu API is now running on port 52983${NC}"
echo -e "${GREEN}📚 Swagger docs will be available at: https://[your-domain]/swagger/index.html${NC}"