# Project Service - Docker Setup

This directory contains the Docker configuration for the Project Service microservice, which provides a marketplace for carbon credit projects.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

### Using Docker Compose (Recommended)

1. **Start all services (Project Service + Dependencies):**
   ```bash
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   docker-compose logs -f project-service
   ```

3. **Stop all services:**
   ```bash
   docker-compose down
   ```

### Using Docker Only

1. **Build the image:**
   ```bash
   docker build -t project-service .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8081:8081 \
     -e DB_HOST=host.docker.internal \
     -e DB_USER=postgres \
     -e DB_PASSWORD=postgres \
     -e DB_NAME=carbon_clear_projects \
     -e DB_PORT=5432 \
     -e ELASTICSEARCH_URL=http://host.docker.internal:9200 \
     -e REDIS_ADDR=host.docker.internal:6379 \
     project-service
   ```

## Services Included

- **Project Service**: Main application (Port 8081)
- **PostgreSQL**: Database (Port 5433)
- **Elasticsearch**: Search engine (Port 9200)
- **Redis**: Caching layer (Port 6379)

## Environment Variables

Copy `env.example` to `.env` and modify as needed:

```bash
cp env.example .env
```

### Required Variables

- `PORT`: Server port (default: 8081)
- `DB_HOST`: PostgreSQL host (default: postgres)
- `DB_USER`: PostgreSQL username (default: postgres)
- `DB_PASSWORD`: PostgreSQL password (default: postgres)
- `DB_NAME`: Database name (default: carbon_clear_projects)
- `DB_PORT`: PostgreSQL port (default: 5432)

### Optional Variables

- `ELASTICSEARCH_URL`: Elasticsearch connection URL
- `ELASTICSEARCH_USERNAME`: Elasticsearch username (if authentication enabled)
- `ELASTICSEARCH_PASSWORD`: Elasticsearch password (if authentication enabled)
- `REDIS_ADDR`: Redis connection address
- `REDIS_PASSWORD`: Redis password (if authentication enabled)
- `REDIS_DB`: Redis database number (default: 0)

## Database Schema

The service uses GORM for database management with the following features:

- **Auto-migration**: Tables are created automatically
- **Project model**: Includes comprehensive project information
- **Search optimization**: Indexes for search performance
- **Timestamps**: Automatic created_at, updated_at tracking
- **Soft deletes**: Support for deleted_at field

### Project Fields

- Basic info: title, description, category
- Location: region, country
- Financial: price_per_tonne, total_capacity, available_capacity
- Verification: verification_standard
- Developer: project_developer, project_url
- Media: image_url
- Status: active/inactive status

## Elasticsearch Integration

The service includes Elasticsearch for advanced search capabilities:

- **Full-text search**: Search across project titles and descriptions
- **Filtered search**: Filter by category, region, country, price range
- **Faceted search**: Get aggregations for categories, regions, etc.
- **Auto-indexing**: Projects are automatically indexed when created/updated

### Elasticsearch Index Structure

```json
{
  "mappings": {
    "properties": {
      "id": {"type": "integer"},
      "title": {"type": "text", "analyzer": "standard"},
      "description": {"type": "text", "analyzer": "standard"},
      "category": {"type": "keyword"},
      "region": {"type": "keyword"},
      "country": {"type": "keyword"},
      "verification_standard": {"type": "keyword"},
      "price_per_tonne": {"type": "float"},
      "total_capacity": {"type": "float"},
      "available_capacity": {"type": "float"},
      "project_developer": {"type": "keyword"},
      "project_url": {"type": "keyword"},
      "image_url": {"type": "keyword"},
      "status": {"type": "keyword"},
      "created_at": {"type": "date"},
      "updated_at": {"type": "date"}
    }
  }
}
```

## Redis Caching

The service uses Redis for caching to improve performance:

- **Project caching**: Cache individual project details
- **Search result caching**: Cache search results with filters
- **List caching**: Cache paginated project lists
- **Metadata caching**: Cache categories, regions, countries

### Cache Keys

- `project:{id}`: Individual project cache
- `projects:all:{limit}:{offset}`: Paginated project list
- `search:{query}:{filters}:{limit}:{offset}`: Search results
- `project:categories`: Available categories
- `project:regions`: Available regions
- `project:countries`: Available countries

## Health Checks

The service includes health checks that verify:
- Project Service: HTTP endpoint at `/health`
- PostgreSQL: Database connectivity and readiness
- Elasticsearch: Search engine connectivity
- Redis: Cache layer connectivity

## Development

### Rebuilding the Service

```bash
docker-compose build project-service
docker-compose up -d project-service
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f project-service
```

### Database Access

Connect to PostgreSQL directly:

```bash
# Using docker-compose
docker-compose exec postgres psql -U postgres -d carbon_clear_projects

# Using docker
docker exec -it project-postgres psql -U postgres -d carbon_clear_projects
```

### Elasticsearch Access

Access Elasticsearch directly:

```bash
# Check cluster health
curl http://localhost:9200/_cluster/health

# List indices
curl http://localhost:9200/_cat/indices

# Search projects
curl "http://localhost:9200/projects/_search?pretty"
```

### Redis Access

Connect to Redis directly:

```bash
# Using docker-compose
docker-compose exec redis redis-cli

# Using docker
docker exec -it project-redis redis-cli
```

## API Endpoints

The service provides the following endpoints:

- `GET /` - Service information
- `GET /health` - Health check endpoint
- Project management endpoints (configured in routes)

## Swagger API Documentation

Interactive API documentation is available via Swagger UI.

### Accessing Swagger UI

When the service is running, access the Swagger documentation at:

```
http://localhost:8081/swagger/index.html
```

### Features

- **Interactive API Testing**: Test all endpoints directly from the browser
- **Request/Response Examples**: View example requests and responses
- **Authentication**: JWT authentication for protected admin endpoints
- **Model Schemas**: Detailed request and response models

### API Categories

- **Projects**: Browse, search, filter, and view carbon offset projects
  - Get all projects (with pagination)
  - Get project by ID
  - Search projects with filters
  - Get categories, regions, and countries
  - Admin: Create, update, and delete projects (requires authentication)

### Regenerating Swagger Documentation

If you make changes to the API handlers, regenerate the Swagger docs:

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g main.go
```

## Data Persistence

- PostgreSQL data is persisted in the `postgres_data` volume
- Elasticsearch data is persisted in the `elasticsearch_data` volume
- Redis data is persisted in the `redis_data` volume
- Database initialization script runs on first startup
- GORM handles table creation and migrations
- Elasticsearch index is created automatically

## Performance Considerations

### Elasticsearch

- Configured with 512MB heap size for development
- Single-node setup for simplicity
- Security disabled for development
- Adjust `ES_JAVA_OPTS` for production workloads

### Redis

- Default configuration suitable for development
- Consider Redis persistence settings for production
- Monitor memory usage and configure eviction policies

### PostgreSQL

- Standard PostgreSQL configuration
- Consider connection pooling for production
- Monitor query performance and add indexes as needed

## Security Features

- **Non-root user**: Container runs as non-privileged user
- **Network isolation**: Services communicate through Docker network
- **Data encryption**: Consider enabling SSL/TLS for production
- **Access control**: Implement proper authentication and authorization

## Troubleshooting

### Service Won't Start

1. Check if ports are already in use:
   ```bash
   lsof -i :8081
   lsof -i :5433
   lsof -i :9200
   lsof -i :6379
   ```

2. Check service logs:
   ```bash
   docker-compose logs project-service
   ```

### Database Connection Issues

1. Ensure PostgreSQL is running and accessible
2. Check the database environment variables
3. Verify network connectivity between containers
4. Check PostgreSQL logs:
   ```bash
   docker-compose logs postgres
   ```

### Elasticsearch Issues

1. Check Elasticsearch logs:
   ```bash
   docker-compose logs elasticsearch
   ```

2. Verify Elasticsearch is accessible:
   ```bash
   curl http://localhost:9200
   ```

3. Check cluster health:
   ```bash
   curl http://localhost:9200/_cluster/health
   ```

### Redis Issues

1. Check Redis logs:
   ```bash
   docker-compose logs redis
   ```

2. Test Redis connection:
   ```bash
   docker-compose exec redis redis-cli ping
   ```

## Production Considerations

### Environment Variables

For production, ensure you set secure values for:

- `DB_PASSWORD`: Use a strong, unique password
- `ELASTICSEARCH_USERNAME/PASSWORD`: Enable authentication
- `REDIS_PASSWORD`: Enable authentication
- `JWT_SECRET`: Use a cryptographically secure secret

### Database Configuration

- Configure PostgreSQL for production workloads
- Set up proper backup strategies
- Monitor database performance
- Configure connection pooling

### Elasticsearch Configuration

- Enable security features
- Configure proper heap size
- Set up cluster for high availability
- Monitor cluster health and performance

### Redis Configuration

- Enable authentication
- Configure persistence settings
- Set up Redis Sentinel for high availability
- Monitor memory usage and performance

### Monitoring

- Set up health checks and monitoring
- Monitor service logs
- Track performance metrics
- Set up alerting for critical issues
