# DramaQu API - Production Deployment Guide

## 🚀 Quick Start

### One-Command Deployment
```bash
./quick-deploy.sh
```

This will automatically:
- Build the Docker image
- Generate Swagger documentation
- Deploy to production on port 52983
- Enable dynamic host detection for wildcard domains

## 📋 Prerequisites

- Docker & Docker Compose installed
- Port 52983 available
- Go 1.24+ (for development)
- Swag CLI (optional, for Swagger generation)

## 🛠️ Available Scripts

### Production Deployment
```bash
./deploy.sh [command]
```

**Commands:**
- `deploy` - Full deployment to production
- `build` - Build Docker image only
- `restart` - Restart application
- `stop` - Stop application
- `logs` - Show application logs
- `status` - Check application status
- `clean` - Clean up Docker resources
- `backup` - Create configuration backup
- `update` - Update and redeploy
- `health` - Health check

### Development
```bash
./dev.sh [command]
```

**Commands:**
- `run` - Start development server
- `docker` - Run with Docker (dev mode)
- `test` - Run tests
- `install` - Install dependencies
- `docs` - Generate Swagger docs
- `format` - Format code
- `lint` - Lint code
- `hot` - Hot reload (requires air)

### Monitoring
```bash
./monitor.sh [command]
```

**Commands:**
- `status` - Quick status check
- `monitor` - Continuous monitoring
- `stats` - Detailed statistics
- `logs` - Show filtered logs
- `perf` - Performance test

## 🌐 Dynamic Domain Support

The API automatically detects the domain from request headers, supporting:

- **Wildcard domains** (*.yourdomain.com)
- **Multiple subdomains** (api1.domain.com, api2.domain.com)
- **HTTP/HTTPS detection** (automatic scheme detection)
- **Reverse proxy support** (Nginx, Traefik, etc.)

### Swagger Documentation
- **URL:** `https://[your-domain]/swagger/index.html`
- **Config:** `https://[your-domain]/swagger-config`
- **Alternative:** `https://[your-domain]/docs/index.html`

## 🐳 Docker Configuration

### Production (docker-compose.prod.yml)
```yaml
services:
  dramaqu-api:
    build: .
    ports:
      - "52983:52983"
    environment:
      - GIN_MODE=release
      - PORT=52983
      - HOST=0.0.0.0
    restart: unless-stopped
```

### Development (docker-compose.yml)
```yaml
services:
  dramaqu-api:
    build: .
    ports:
      - "52983:52983"
    environment:
      - GIN_MODE=debug
      - PORT=52983
```

## 🔧 Environment Configuration

### Production (.env.production)
```bash
PORT=52983
HOST=0.0.0.0
GIN_MODE=release
TZ=Asia/Jakarta
```

### Development (.env.development)
```bash
PORT=52983
HOST=localhost
GIN_MODE=debug
TZ=Asia/Jakarta
```

## 🌍 Reverse Proxy Setup

### Nginx Configuration
```nginx
server {
    listen 80;
    server_name *.yourdomain.com;
    
    location / {
        proxy_pass http://localhost:52983;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Traefik Configuration (docker-compose.prod.yml)
```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.dramaqu-api.rule=HostRegexp(`{host:.+}`)"
  - "traefik.http.routers.dramaqu-api.entrypoints=websecure"
  - "traefik.http.routers.dramaqu-api.tls.certresolver=letsencrypt"
```

## 📊 Monitoring & Health Checks

### Health Check Endpoint
```bash
curl http://localhost:52983/swagger-config
```

### Container Health Check
```bash
docker exec dramaqu-api-prod wget --spider http://localhost:52983/swagger-config
```

### Continuous Monitoring
```bash
./monitor.sh monitor
```

## 🔄 Deployment Workflow

### Initial Deployment
1. Clone repository
2. Configure environment variables
3. Run: `./quick-deploy.sh`
4. Configure reverse proxy
5. Set up SSL certificate

### Updates
1. Pull latest code: `git pull`
2. Update: `./deploy.sh update`
3. Verify: `./monitor.sh status`

### Rollback
1. Stop current: `./deploy.sh stop`
2. Restore backup configuration
3. Deploy previous version: `./deploy.sh deploy`

## 🚨 Troubleshooting

### Common Issues

**Port already in use:**
```bash
sudo lsof -i :52983
./deploy.sh stop
```

**Container won't start:**
```bash
./deploy.sh logs
docker system prune -f
./deploy.sh clean
./deploy.sh deploy
```

**Swagger not loading:**
```bash
# Check if docs are generated
ls -la docs/
./dev.sh docs
./deploy.sh restart
```

**Health check failing:**
```bash
./monitor.sh status
curl -v http://localhost:52983/swagger-config
./deploy.sh logs
```

### Log Locations
- **Application logs:** `./deploy.sh logs`
- **Monitor logs:** `/tmp/dramaqu-monitor.log`
- **Docker logs:** `docker logs dramaqu-api-prod`

## 📈 Performance Tuning

### Docker Resource Limits
```yaml
services:
  dramaqu-api:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### Go Application Tuning
```bash
# Environment variables
GOMAXPROCS=2
GOGC=100
```

## 🔐 Security Considerations

1. **Run as non-root user** in production
2. **Use HTTPS** with valid SSL certificates
3. **Configure firewall** to allow only necessary ports
4. **Regular updates** of base images and dependencies
5. **Monitor logs** for suspicious activity
6. **Use secrets management** for sensitive data

## 📞 Support

For issues and questions:
1. Check logs: `./deploy.sh logs`
2. Run diagnostics: `./monitor.sh stats`
3. Review this documentation
4. Check Docker and Go documentation

## 🎯 API Endpoints

Once deployed, the following endpoints are available:

- **Swagger UI:** `/swagger/index.html`
- **API Documentation:** `/docs/index.html`
- **Swagger Config:** `/swagger-config`
- **Health Check:** `/swagger-config` (returns JSON)

All endpoints support dynamic host detection and will work with any domain pointing to your server.