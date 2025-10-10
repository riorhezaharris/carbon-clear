# Carbon Clear API Gateway

The API Gateway serves as the single entry point for all client-side communications with the Carbon Clear microservices architecture. It handles routing, authentication, rate limiting, and provides a unified API interface.

## Features

- **Unified API Interface**: Single entry point for all microservices
- **Authentication & Authorization**: JWT-based authentication for users and admins
- **Rate Limiting**: Configurable rate limiting to prevent abuse
- **Request Routing**: Intelligent routing to appropriate backend services
- **CORS Support**: Cross-Origin Resource Sharing enabled
- **Health Checks**: Monitoring endpoint for service health
- **Logging**: Comprehensive request/response logging
- **Swagger Documentation**: Interactive API documentation

## Architecture

The API Gateway routes requests to three backend services:

1. **User Service** (port 8082): User authentication and management
2. **Project Service** (port 8081): Carbon offset project marketplace
3. **Order Service** (port 8080): Shopping cart and order management

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)
- Access to backend services (User, Project, Order services)

## Installation

### Local Development

1. Clone the repository and navigate to the api_gateway directory:
```bash
cd api_gateway
```

2. Copy the example environment file and configure it:
```bash
cp env.example .env
```

3. Edit `.env` file with your configuration:
```env
PORT=8000
USER_SERVICE_URL=http://localhost:8082
PROJECT_SERVICE_URL=http://localhost:8081
ORDER_SERVICE_URL=http://localhost:8080
USER_JWT_SECRET=your-user-jwt-secret
ADMIN_JWT_SECRET=your-admin-jwt-secret
RATE_LIMIT_PER_MIN=100
```

4. Install dependencies:
```bash
go mod download
```

5. Run the gateway:
```bash
go run main.go
```

### Docker Deployment

1. Build the Docker image:
```bash
docker build -t carbon-clear-api-gateway .
```

2. Run with Docker Compose:
```bash
docker-compose up -d
```

## API Endpoints

### Health Check
- `GET /health` - Check gateway health status
- `GET /` - Gateway information

### User Service Routes

#### Public Routes
- `POST /api/users/register` - Register new user
- `POST /api/users/login` - User login

#### Protected Routes (User JWT Required)
- `GET /api/users/profile` - Get user profile

#### Admin Routes (Admin JWT Required)
- `POST /api/admin/users/register` - Register new admin
- `POST /api/admin/users/login` - Admin login
- `GET /api/admin/users` - Get all users
- `POST /api/admin/users` - Create new user
- `GET /api/admin/users/:id` - Get user by ID
- `PUT /api/admin/users/:id` - Update user
- `DELETE /api/admin/users/:id` - Delete user

### Project Service Routes

#### Public Routes
- `GET /api/v1/projects` - Get all projects
- `GET /api/v1/projects/:id` - Get project details
- `POST /api/v1/projects/search` - Search projects
- `GET /api/v1/projects/categories` - Get project categories
- `GET /api/v1/projects/regions` - Get project regions
- `GET /api/v1/projects/countries` - Get project countries

#### Admin Routes (Admin JWT Required)
- `POST /api/v1/projects/admin` - Create new project
- `PUT /api/v1/projects/admin/:id` - Update project
- `DELETE /api/v1/projects/admin/:id` - Delete project

### Order Service Routes

#### Cart Routes (User JWT Required)
- `POST /api/v1/cart/:userID/items` - Add item to cart
- `GET /api/v1/cart/:userID` - Get user's cart
- `PUT /api/v1/cart/:userID/items/:projectID` - Update cart item
- `DELETE /api/v1/cart/:userID/items/:projectID` - Remove item from cart
- `DELETE /api/v1/cart/:userID` - Clear cart

#### Order Routes (User JWT Required)
- `POST /api/v1/orders/:userID/checkout` - Checkout cart
- `GET /api/v1/orders/:userID/history` - Get order history
- `GET /api/v1/orders/:orderID` - Get order details
- `GET /api/v1/orders/:userID/certificates` - Get certificates

#### Admin Order Routes (Admin JWT Required)
- `GET /api/v1/admin/reports/monthly` - Get monthly reports
- `GET /api/v1/admin/orders/date-range` - Get orders by date range
- `GET /api/v1/admin/statistics` - Get order statistics

## Authentication

The gateway uses JWT (JSON Web Tokens) for authentication. There are two types of tokens:

1. **User Token**: For regular user operations
2. **Admin Token**: For administrative operations

### Using Authentication

Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Rate Limiting

The gateway implements rate limiting to prevent abuse:
- Default: 100 requests per minute per IP
- Configurable via `RATE_LIMIT_PER_MIN` environment variable
- Returns `429 Too Many Requests` when limit exceeded

## CORS Configuration

CORS is enabled for all origins. In production, configure specific origins in `main.go`:

```go
AllowOrigins: []string{"https://yourdomain.com"},
```

## Monitoring

### Health Check

Check gateway status:
```bash
curl http://localhost:8000/health
```

Response:
```json
{
  "status": "healthy",
  "service": "api_gateway",
  "version": "1.0.0"
}
```

### Logs

The gateway logs all requests with details:
- Timestamp
- HTTP method
- Path
- Status code
- Latency
- Client IP
- User agent

## Swagger Documentation

Access interactive API documentation at:
```
http://localhost:8000/swagger/index.html
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Gateway port | 8000 |
| USER_SERVICE_URL | User service URL | http://localhost:8082 |
| PROJECT_SERVICE_URL | Project service URL | http://localhost:8081 |
| ORDER_SERVICE_URL | Order service URL | http://localhost:8080 |
| USER_JWT_SECRET | JWT secret for user tokens | user-secret-key |
| ADMIN_JWT_SECRET | JWT secret for admin tokens | admin-secret-key |
| RATE_LIMIT_PER_MIN | Rate limit per minute | 100 |

## Error Handling

The gateway returns standard HTTP error codes:

- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `429 Too Many Requests` - Rate limit exceeded
- `502 Bad Gateway` - Backend service unavailable
- `503 Service Unavailable` - Gateway overloaded

## Development

### Project Structure

```
api_gateway/
├── config/           # Configuration management
├── middleware/       # Authentication, rate limiting, logging
├── proxy/           # Request proxying and routing
├── docs/            # Swagger documentation
├── main.go          # Application entry point
├── Dockerfile       # Docker configuration
├── docker-compose.yml
├── .env.example
└── README.md
```

### Adding New Routes

1. Add route in `proxy/routes.go`
2. Configure appropriate middleware (auth, rate limiting)
3. Map to backend service URL

Example:
```go
api.GET("/v1/new-endpoint", proxy.ProxyRequest(cfg.ServiceURL))
```

## Security Best Practices

1. **Never commit JWT secrets** - Use environment variables
2. **Use HTTPS in production** - Configure TLS/SSL certificates
3. **Implement proper CORS** - Restrict allowed origins
4. **Monitor rate limits** - Adjust based on traffic patterns
5. **Regular security updates** - Keep dependencies updated
6. **Log monitoring** - Track suspicious activities

## Troubleshooting

### Gateway can't connect to backend service

1. Check if backend services are running
2. Verify service URLs in `.env`
3. Check network connectivity
4. Review firewall rules

### Authentication errors

1. Verify JWT secrets match between gateway and services
2. Check token expiration
3. Ensure correct token format: `Bearer <token>`

### Rate limiting too strict

Adjust `RATE_LIMIT_PER_MIN` in `.env`:
```env
RATE_LIMIT_PER_MIN=200
```

## Performance

The gateway is designed for high performance:
- **Request timeout**: 30 seconds
- **Connection pooling**: Managed by Go http.Client
- **Concurrent requests**: Limited by rate limiter
- **Memory usage**: ~50MB base + request overhead

## Contributing

1. Follow Go best practices
2. Add tests for new features
3. Update documentation
4. Test with all backend services

## License

Apache 2.0

## Support

For support, contact: support@carbonclear.com

