# Carbon Clear - Deployment Guide

This guide provides step-by-step instructions for deploying the Carbon Clear microservices platform.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start](#quick-start)
3. [Configuration](#configuration)
4. [Deployment Options](#deployment-options)
5. [Production Deployment](#production-deployment)
6. [Monitoring and Maintenance](#monitoring-and-maintenance)
7. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher
- **Git**: For cloning the repository
- **Make**: (Optional) For using Makefile commands

### System Requirements

**Minimum:**
- CPU: 2 cores
- RAM: 4 GB
- Disk: 10 GB free space

**Recommended:**
- CPU: 4 cores
- RAM: 8 GB
- Disk: 20 GB free space

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd carbon-clear
```

### 2. Create Environment Configuration

Create a `.env` file in the project root:

```bash
cat > .env << EOF
# JWT Secrets - CHANGE THESE IN PRODUCTION
USER_JWT_SECRET=$(openssl rand -base64 32)
ADMIN_JWT_SECRET=$(openssl rand -base64 32)

# Database Passwords - CHANGE THESE IN PRODUCTION
POSTGRES_PASSWORD=postgres
MONGODB_PASSWORD=admin123
REDIS_PASSWORD=redis123

# RabbitMQ Configuration
RABBITMQ_USER=admin
RABBITMQ_PASSWORD=admin123
EOF
```

### 3. Start All Services

Using Make (recommended):
```bash
make build
make up
```

Or using Docker Compose directly:
```bash
docker-compose build
docker-compose up -d
```

### 4. Verify Deployment

Check service status:
```bash
make status
# or
docker-compose ps
```

Run health checks:
```bash
make test
```

### 5. Access Services

- **API Gateway**: http://localhost:8000
- **API Documentation**: http://localhost:8000/swagger/index.html
- **RabbitMQ Console**: http://localhost:15672 (admin/admin123)

## Configuration

### Environment Variables

#### JWT Configuration (Required)

```env
USER_JWT_SECRET=your-strong-random-secret-here
ADMIN_JWT_SECRET=your-strong-random-secret-here
```

⚠️ **Important**: Generate strong random secrets:
```bash
openssl rand -base64 32
```

#### Database Configuration

PostgreSQL (User & Project Services):
```env
POSTGRES_PASSWORD=secure-password
```

MongoDB (Order Service):
```env
MONGODB_PASSWORD=secure-password
```

Redis (Project Service Cache):
```env
REDIS_PASSWORD=secure-password
```

#### Service URLs

API Gateway service URLs are configured in docker-compose.yml:
```yaml
USER_SERVICE_URL=http://user_service:8082
PROJECT_SERVICE_URL=http://project_service:8081
ORDER_SERVICE_URL=http://order_service:8080
```

#### Rate Limiting

Configure API Gateway rate limiting:
```env
RATE_LIMIT_PER_MIN=100
```

### Port Configuration

Default ports can be changed in `docker-compose.yml`:

```yaml
services:
  api_gateway:
    ports:
      - "8000:8000"  # Change left number for external port
```

## Deployment Options

### Option 1: Full Stack (All Services)

Deploy all services together:

```bash
docker-compose up -d
```

### Option 2: Selective Services

Deploy specific services:

```bash
# Start databases first
docker-compose up -d postgres_user postgres_project mongodb redis rabbitmq elasticsearch

# Start individual services
docker-compose up -d user_service
docker-compose up -d project_service
docker-compose up -d order_service
docker-compose up -d api_gateway
```

### Option 3: Development Mode

Run services locally for development:

```bash
# Start only databases
docker-compose up -d postgres_user postgres_project mongodb redis rabbitmq elasticsearch

# Run services locally
make dev-user        # Terminal 1
make dev-project     # Terminal 2
make dev-order       # Terminal 3
make dev-gateway     # Terminal 4
```

## Production Deployment

### 1. Security Configuration

#### Generate Strong Secrets

```bash
# Generate JWT secrets
export USER_JWT_SECRET=$(openssl rand -base64 32)
export ADMIN_JWT_SECRET=$(openssl rand -base64 32)

# Generate database passwords
export POSTGRES_PASSWORD=$(openssl rand -base64 24)
export MONGODB_PASSWORD=$(openssl rand -base64 24)
export REDIS_PASSWORD=$(openssl rand -base64 24)
export RABBITMQ_PASSWORD=$(openssl rand -base64 24)
```

#### Update CORS Configuration

Edit `api_gateway/main.go`:

```go
e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
    AllowOrigins: []string{"https://yourdomain.com"},  // Change this
    AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
    AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
}))
```

### 2. Database Persistence

Ensure data persistence by using volumes:

```yaml
volumes:
  postgres_user_data:
    driver: local
  postgres_project_data:
    driver: local
  mongo_order_data:
    driver: local
```

For production, consider using external database services.

### 3. SSL/TLS Configuration

#### Option A: Using Reverse Proxy (Recommended)

Use nginx or traefik as a reverse proxy:

```nginx
server {
    listen 443 ssl;
    server_name api.yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### Option B: Direct TLS in Gateway

Configure Echo with TLS:

```go
e.Logger.Fatal(e.StartTLS(":443", "cert.pem", "key.pem"))
```

### 4. Resource Limits

Add resource limits to docker-compose.yml:

```yaml
services:
  api_gateway:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 5. Logging Configuration

#### Centralized Logging

Add logging driver to docker-compose.yml:

```yaml
services:
  api_gateway:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

#### ELK Stack Integration

For production, consider integrating ELK stack for centralized logging.

### 6. Backup Strategy

#### Database Backups

PostgreSQL:
```bash
# Automated backup script
docker exec postgres_user pg_dump -U postgres user_service_db > backup_user_$(date +%Y%m%d).sql
docker exec postgres_project pg_dump -U postgres project_service_db > backup_project_$(date +%Y%m%d).sql
```

MongoDB:
```bash
docker exec mongodb mongodump --username admin --password admin123 --authenticationDatabase admin --out /backup
```

Set up automated backups with cron:
```bash
0 2 * * * /path/to/backup-script.sh
```

### 7. Monitoring Setup

#### Prometheus + Grafana

Add monitoring services to docker-compose.yml:

```yaml
  prometheus:
    image: prom/prometheus
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

## Monitoring and Maintenance

### Health Checks

```bash
# Check all services
make status

# Run health checks
make test

# Individual service health
curl http://localhost:8000/health
curl http://localhost:8082/
curl http://localhost:8081/health
```

### View Logs

```bash
# All services
make logs

# Individual service
make gateway-logs
make user-logs
make project-logs
make order-logs
```

### Service Management

```bash
# Restart all services
make restart

# Restart individual service
make gateway-restart
make user-restart
make project-restart
make order-restart
```

### Database Management

```bash
# Access PostgreSQL (User Service)
docker exec -it postgres_user psql -U postgres -d user_service_db

# Access PostgreSQL (Project Service)
docker exec -it postgres_project psql -U postgres -d project_service_db

# Access MongoDB
docker exec -it mongodb mongosh -u admin -p admin123 --authenticationDatabase admin

# Access Redis
docker exec -it redis redis-cli -a redis123
```

### RabbitMQ Management

Access RabbitMQ console: http://localhost:15672
- Username: admin
- Password: admin123

Monitor queues and messages for certificate generation.

## Troubleshooting

### Services Won't Start

1. **Check Docker resources:**
   ```bash
   docker system df
   docker system prune
   ```

2. **Check logs:**
   ```bash
   docker-compose logs
   ```

3. **Verify environment variables:**
   ```bash
   docker-compose config
   ```

### Database Connection Errors

1. **Check database health:**
   ```bash
   docker-compose ps
   ```

2. **Verify database is ready:**
   ```bash
   docker exec postgres_user pg_isready -U postgres
   ```

3. **Check connection strings:**
   Review database URLs in service environment variables.

### Authentication Errors

1. **Verify JWT secrets match:**
   Ensure the same secrets are used across gateway and services.

2. **Check token format:**
   ```bash
   # Token should be in format: Bearer <token>
   curl -H "Authorization: Bearer your-token-here" http://localhost:8000/api/users/profile
   ```

3. **Token expiration:**
   Login again to get a fresh token.

### Rate Limiting Issues

If getting 429 errors:

1. **Adjust rate limit:**
   ```env
   RATE_LIMIT_PER_MIN=200
   ```

2. **Restart gateway:**
   ```bash
   make gateway-restart
   ```

### Performance Issues

1. **Check resource usage:**
   ```bash
   docker stats
   ```

2. **Increase resources:**
   Edit docker-compose.yml resource limits.

3. **Scale services:**
   ```bash
   docker-compose up -d --scale api_gateway=3
   ```

### Port Conflicts

If ports are in use:

1. **Find processes:**
   ```bash
   lsof -i :8000
   lsof -i :8080
   ```

2. **Change ports in docker-compose.yml:**
   ```yaml
   ports:
     - "8001:8000"  # Use different external port
   ```

## Production Checklist

Before deploying to production:

- [ ] Generate strong JWT secrets
- [ ] Change all default passwords
- [ ] Configure CORS for specific domain
- [ ] Set up SSL/TLS certificates
- [ ] Configure resource limits
- [ ] Set up automated backups
- [ ] Configure monitoring and alerting
- [ ] Set up centralized logging
- [ ] Test health checks
- [ ] Document recovery procedures
- [ ] Plan scaling strategy
- [ ] Security audit
- [ ] Load testing

## Scaling Considerations

### Horizontal Scaling

```bash
# Scale API Gateway
docker-compose up -d --scale api_gateway=3

# Scale services
docker-compose up -d --scale user_service=2
docker-compose up -d --scale project_service=2
docker-compose up -d --scale order_service=2
```

### Load Balancing

Use nginx or HAProxy for load balancing:

```nginx
upstream api_gateway {
    server localhost:8000;
    server localhost:8001;
    server localhost:8002;
}
```

### Database Scaling

Consider:
- Read replicas for PostgreSQL
- MongoDB replica sets
- Redis cluster mode
- Elasticsearch cluster

## Disaster Recovery

### Backup Strategy

1. **Daily automated backups**
2. **Weekly full backups**
3. **Monthly archival**
4. **Off-site backup storage**

### Recovery Procedures

1. **Stop services:**
   ```bash
   make down
   ```

2. **Restore databases:**
   ```bash
   # PostgreSQL
   cat backup.sql | docker exec -i postgres_user psql -U postgres user_service_db
   
   # MongoDB
   docker exec mongodb mongorestore --username admin --password admin123 /backup
   ```

3. **Restart services:**
   ```bash
   make up
   ```

4. **Verify:**
   ```bash
   make test
   ```

## Support

For additional support:
- Email: support@carbonclear.com
- Documentation: See README.md
- Issues: GitHub Issues

## Updates and Upgrades

### Updating Services

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
make down
make build
make up
```

### Database Migrations

Follow service-specific migration procedures before updating.

## Conclusion

This deployment guide covers the essential steps for deploying Carbon Clear. Always test in a staging environment before deploying to production.

For development documentation, see [README.md](README.md).

