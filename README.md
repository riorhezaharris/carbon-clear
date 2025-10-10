# Carbon Clear - Microservices Platform

A comprehensive carbon offset marketplace platform built with microservices architecture, featuring user management, project marketplace, and order processing capabilities.

## 🏗️ Architecture Overview

The platform consists of four main components:

1. **API Gateway** (Port 8000) - Single entry point for all client requests
2. **User Service** (Port 8082) - User authentication and profile management
3. **Project Service** (Port 8081) - Carbon offset project marketplace
4. **Order Service** (Port 8080) - Shopping cart and order processing

### Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Echo v4
- **Databases**: PostgreSQL, MongoDB
- **Message Queue**: RabbitMQ
- **Cache**: Redis
- **Search**: Elasticsearch
- **API Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Make (optional)

### Running with Docker Compose

1. Clone the repository:
```bash
git clone <repository-url>
cd carbon-clear
```

2. Create environment file:
```bash
# Create .env file with your configurations
cat > .env << EOF
USER_JWT_SECRET=your-strong-user-secret-key
ADMIN_JWT_SECRET=your-strong-admin-secret-key
EOF
```

3. Start all services:
```bash
docker-compose up -d
```

4. Verify all services are running:
```bash
docker-compose ps
```

5. Access the services:
- **API Gateway**: http://localhost:8000
- **API Gateway Swagger**: http://localhost:8000/swagger/index.html
- **User Service**: http://localhost:8082
- **Project Service**: http://localhost:8081
- **Order Service**: http://localhost:8080
- **RabbitMQ Management**: http://localhost:15672 (admin/admin123)

### Service Health Checks

```bash
# API Gateway
curl http://localhost:8000/health

# User Service
curl http://localhost:8082/

# Project Service  
curl http://localhost:8081/health

# Order Service
curl http://localhost:8080/swagger/
```

## 📋 Services Details

### API Gateway

The API Gateway is the single entry point for all client communications.

**Features:**
- Request routing to backend services
- JWT authentication (user and admin)
- Rate limiting (100 req/min by default)
- CORS support
- Request/response logging
- Health monitoring

**Documentation:** See [api_gateway/README.md](api_gateway/README.md)

### User Service

Manages user authentication and profile information.

**Features:**
- User registration and login
- Admin registration and login
- JWT token generation
- Profile management
- User CRUD operations (admin only)

**Key Endpoints:**
- `POST /api/users/register` - Register user
- `POST /api/users/login` - User login
- `GET /api/users/profile` - Get profile (authenticated)
- `POST /admin/users/register` - Register admin
- `POST /admin/users/login` - Admin login

**Documentation:** See [user_service/README.md](user_service/README.md)

### Project Service

Handles carbon offset project marketplace.

**Features:**
- Browse carbon offset projects
- Search and filter projects
- Project categories, regions, countries
- Admin project management (CRUD)
- Redis caching
- Elasticsearch integration

**Key Endpoints:**
- `GET /api/v1/projects` - List all projects
- `GET /api/v1/projects/:id` - Get project details
- `POST /api/v1/projects/search` - Search projects
- `POST /api/v1/projects/admin` - Create project (admin)

**Documentation:** See [project_service/README.md](project_service/README.md)

### Order Service

Manages shopping cart, orders, and certificates.

**Features:**
- Shopping cart management
- Order checkout and processing
- Order history
- Certificate generation (via RabbitMQ)
- Monthly reports and statistics
- Scheduled tasks

**Key Endpoints:**
- `POST /api/v1/cart/:userID/items` - Add to cart
- `GET /api/v1/cart/:userID` - Get cart
- `POST /api/v1/orders/:userID/checkout` - Checkout
- `GET /api/v1/orders/:userID/history` - Order history
- `GET /api/v1/orders/:userID/certificates` - Get certificates

**Documentation:** See [order_service/README.md](order_service/README.md)

## 🔐 Authentication

The platform uses JWT (JSON Web Tokens) for authentication with two levels:

### User Authentication
For regular user operations (cart, orders, profile).

**Login:**
```bash
curl -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

**Use Token:**
```bash
curl -X GET http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer <user-token>"
```

### Admin Authentication
For administrative operations (manage users, projects, view reports).

**Login:**
```bash
curl -X POST http://localhost:8000/api/admin/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "admin123"}'
```

**Use Token:**
```bash
curl -X GET http://localhost:8000/api/admin/users \
  -H "Authorization: Bearer <admin-token>"
```

## 🛠️ Development

### Local Development Setup

Each service can be run independently for development:

#### API Gateway
```bash
cd api_gateway
cp env.example .env
# Edit .env with your configuration
go mod download
go run main.go
```

#### User Service
```bash
cd user_service
cp env.example .env
# Start PostgreSQL first
go mod download
go run main.go
```

#### Project Service
```bash
cd project_service
cp env.example .env
# Start PostgreSQL, Redis, Elasticsearch
go mod download
go run main.go
```

#### Order Service
```bash
cd order_service
cp env.example .env
# Start MongoDB, RabbitMQ
go mod download
go run main.go
```

### Project Structure

```
carbon-clear/
├── api_gateway/          # API Gateway service
│   ├── config/           # Configuration management
│   ├── middleware/       # Auth, rate limiting, logging
│   ├── proxy/           # Request routing
│   └── main.go
├── user_service/         # User management service
│   ├── configs/          # Database config
│   ├── handlers/         # HTTP handlers
│   ├── models/          # Data models
│   ├── repositories/    # Data access layer
│   └── routes/          # Route definitions
├── project_service/      # Project marketplace service
│   ├── config/          # Database, Redis, ES config
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── repositories/    # Data access layer
│   └── routes/          # Route definitions
├── order_service/        # Order processing service
│   ├── config/          # Database, RabbitMQ config
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── repositories/    # Data access layer
│   ├── services/        # Business logic
│   └── routes/          # Route definitions
└── docker-compose.yml    # Orchestration for all services
```

## 📊 Database Schema

### User Service (PostgreSQL)
- `users` table: id, name, email, password, role, created_at, updated_at

### Project Service (PostgreSQL)
- `projects` table: id, name, description, category, location, region, country, price_per_ton, available_tons, image_url, status, created_at, updated_at

### Order Service (MongoDB)
- `carts` collection: user cart items
- `orders` collection: completed orders
- `certificates` collection: generated certificates

## 🔧 Configuration

### Environment Variables

All services support configuration via environment variables. See individual service README files for detailed configuration options.

**Common Variables:**
- `PORT` - Service port
- `USER_JWT_SECRET` - JWT secret for user tokens
- `ADMIN_JWT_SECRET` - JWT secret for admin tokens

**Database Variables:**
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` (PostgreSQL services)
- `MONGODB_URI`, `MONGODB_DATABASE` (Order service)

**External Services:**
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` (Project service)
- `ELASTICSEARCH_URL` (Project service)
- `RABBITMQ_URL` (Order service)

## 🧪 Testing

### Manual API Testing

Use the Swagger UI for interactive testing:
- API Gateway: http://localhost:8000/swagger/index.html
- User Service: http://localhost:8082/swagger/index.html
- Project Service: http://localhost:8081/swagger/index.html
- Order Service: http://localhost:8080/swagger/index.html

### Example User Flow

1. **Register a user:**
```bash
curl -X POST http://localhost:8000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

2. **Login:**
```bash
TOKEN=$(curl -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "password123"}' \
  | jq -r '.token')
```

3. **Browse projects:**
```bash
curl http://localhost:8000/api/v1/projects
```

4. **Add to cart:**
```bash
curl -X POST http://localhost:8000/api/v1/cart/1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": 1,
    "tons": 10
  }'
```

5. **Checkout:**
```bash
curl -X POST http://localhost:8000/api/v1/orders/1/checkout \
  -H "Authorization: Bearer $TOKEN"
```

## 📈 Monitoring

### Health Checks

All services expose health check endpoints:

```bash
# Check all services
docker-compose ps

# Individual health checks
curl http://localhost:8000/health  # API Gateway
curl http://localhost:8082/        # User Service
curl http://localhost:8081/health  # Project Service
```

### Logs

View logs for all services:
```bash
docker-compose logs -f
```

View logs for specific service:
```bash
docker-compose logs -f api_gateway
docker-compose logs -f user_service
docker-compose logs -f project_service
docker-compose logs -f order_service
```

### RabbitMQ Management

Access RabbitMQ management console:
- URL: http://localhost:15672
- Username: admin
- Password: admin123

## 🚢 Deployment

### Production Considerations

1. **Security:**
   - Use strong JWT secrets
   - Enable HTTPS/TLS
   - Restrict CORS origins
   - Use secure database passwords
   - Enable authentication for all external services

2. **Scaling:**
   - Use container orchestration (Kubernetes)
   - Implement horizontal pod autoscaling
   - Use managed database services
   - Set up load balancing

3. **Monitoring:**
   - Implement centralized logging (ELK stack)
   - Set up metrics (Prometheus + Grafana)
   - Configure alerting
   - Use distributed tracing

4. **Backup:**
   - Regular database backups
   - Configuration backups
   - Disaster recovery plan

## 🐛 Troubleshooting

### Services won't start

```bash
# Check logs
docker-compose logs

# Restart specific service
docker-compose restart api_gateway

# Rebuild and restart
docker-compose up -d --build
```

### Database connection errors

```bash
# Check database health
docker-compose ps

# Restart database
docker-compose restart postgres_user
docker-compose restart mongodb
```

### Port conflicts

If ports are already in use, modify the port mappings in `docker-compose.yml`:
```yaml
ports:
  - "8001:8000"  # Change external port
```

## 📝 API Documentation

Complete API documentation is available via Swagger UI:

- **API Gateway**: http://localhost:8000/swagger/index.html
- **User Service**: http://localhost:8082/swagger/index.html
- **Project Service**: http://localhost:8081/swagger/index.html
- **Order Service**: http://localhost:8080/swagger/index.html

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

Apache 2.0

## 📧 Support

For support, please contact: support@carbonclear.com

