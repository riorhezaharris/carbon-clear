# User Service - Docker Setup

This directory contains the Docker configuration for the User Service microservice.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

### Using Docker Compose (Recommended)

1. **Start all services (User Service + PostgreSQL):**
   ```bash
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   docker-compose logs -f user-service
   ```

3. **Stop all services:**
   ```bash
   docker-compose down
   ```

### Using Docker Only

1. **Build the image:**
   ```bash
   docker build -t user-service .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 \
     -e DB_HOST=host.docker.internal \
     -e DB_USER=postgres \
     -e DB_PASSWORD=postgres \
     -e DB_NAME=carbon_clear_users \
     -e DB_PORT=5432 \
     user-service
   ```

## Services Included

- **User Service**: Main application (Port 8080)
- **PostgreSQL**: Database (Port 5432)

## Environment Variables

Copy `env.example` to `.env` and modify as needed:

```bash
cp env.example .env
```

### Required Variables

- `PORT`: Server port (default: 8080)
- `DB_HOST`: PostgreSQL host (default: postgres)
- `DB_USER`: PostgreSQL username (default: postgres)
- `DB_PASSWORD`: PostgreSQL password (default: postgres)
- `DB_NAME`: Database name (default: carbon_clear_users)
- `DB_PORT`: PostgreSQL port (default: 5432)

### Optional Variables

- `JWT_SECRET`: JWT secret key for authentication
- `JWT_EXPIRES_IN`: JWT token expiration time

## Database Schema

The service uses GORM for database management with the following features:

- **Auto-migration**: Tables are created automatically
- **User model**: Includes email, password, profile fields
- **Timestamps**: Automatic created_at, updated_at tracking
- **Soft deletes**: Support for deleted_at field
- **Indexes**: Optimized for email lookups and queries

## Health Checks

The service includes health checks that verify:
- User Service: HTTP endpoint at `/`
- PostgreSQL: Database connectivity and readiness

## Development

### Rebuilding the Service

```bash
docker-compose build user-service
docker-compose up -d user-service
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f user-service
```

### Database Access

Connect to PostgreSQL directly:

```bash
# Using docker-compose
docker-compose exec postgres psql -U postgres -d carbon_clear_users

# Using docker
docker exec -it user-postgres psql -U postgres -d carbon_clear_users
```

### Database Management

```bash
# View database status
docker-compose exec postgres pg_isready -U postgres

# Backup database
docker-compose exec postgres pg_dump -U postgres carbon_clear_users > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres carbon_clear_users < backup.sql
```

## API Endpoints

The service provides the following endpoints:

- `GET /` - Health check endpoint
- User management endpoints (configured in routes)

## Data Persistence

- PostgreSQL data is persisted in the `postgres_data` volume
- Database initialization script runs on first startup
- GORM handles table creation and migrations

## Security Features

- **Non-root user**: Container runs as non-privileged user
- **Password hashing**: Secure password storage
- **JWT authentication**: Token-based authentication
- **Database security**: Proper user permissions and access control

## Troubleshooting

### Service Won't Start

1. Check if ports are already in use:
   ```bash
   lsof -i :8080
   lsof -i :5432
   ```

2. Check service logs:
   ```bash
   docker-compose logs user-service
   ```

### Database Connection Issues

1. Ensure PostgreSQL is running and accessible
2. Check the database environment variables
3. Verify network connectivity between containers
4. Check PostgreSQL logs:
   ```bash
   docker-compose logs postgres
   ```

### Migration Issues

1. Check if GORM can connect to the database
2. Verify database permissions
3. Check for schema conflicts
4. Review application logs for migration errors

## Production Considerations

### Environment Variables

For production, ensure you set secure values for:

- `DB_PASSWORD`: Use a strong, unique password
- `JWT_SECRET`: Use a cryptographically secure secret
- `DB_HOST`: Use your production database host

### Database Configuration

- Configure PostgreSQL for production workloads
- Set up proper backup strategies
- Monitor database performance
- Configure connection pooling if needed

### Security

- Use secrets management for sensitive data
- Enable SSL/TLS for database connections
- Implement proper authentication and authorization
- Regular security updates and monitoring
