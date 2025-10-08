# Order Service - Docker Setup

This directory contains the Docker configuration for the Order Service microservice.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

### Using Docker Compose (Recommended)

1. **Start all services (Order Service + Dependencies):**
   ```bash
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   docker-compose logs -f order-service
   ```

3. **Stop all services:**
   ```bash
   docker-compose down
   ```

### Using Docker Only

1. **Build the image:**
   ```bash
   docker build -t order-service .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 \
     -e MONGODB_URI=mongodb://host.docker.internal:27017 \
     -e RABBITMQ_URL=amqp://guest:guest@host.docker.internal:5672/ \
     order-service
   ```

## Services Included

- **Order Service**: Main application (Port 8080)
- **MongoDB**: Database (Port 27017)
- **RabbitMQ**: Message broker (Port 5672, Management UI: 15672)

## Environment Variables

Copy `env.example` to `.env` and modify as needed:

```bash
cp env.example .env
```

### Required Variables

- `PORT`: Server port (default: 8080)
- `MONGODB_URI`: MongoDB connection string
- `RABBITMQ_URL`: RabbitMQ connection string

## Health Checks

The service includes health checks that verify:
- Order Service: HTTP endpoint at `/health`
- MongoDB: Database connectivity
- RabbitMQ: Message broker connectivity

## Development

### Rebuilding the Service

```bash
docker-compose build order-service
docker-compose up -d order-service
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f order-service
```

### Accessing Services

- **Order Service API**: http://localhost:8080
- **MongoDB**: localhost:27017
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

## Swagger API Documentation

Interactive API documentation is available via Swagger UI.

### Accessing Swagger UI

When the service is running, access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

### Features

- **Interactive API Testing**: Test all endpoints directly from the browser
- **Request/Response Examples**: View example requests and responses
- **Authentication**: JWT authentication for protected admin endpoints
- **Model Schemas**: Detailed request and response models

### API Categories

- **Cart**: Shopping cart management
  - Add items to cart
  - Get cart items
  - Update cart item quantity
  - Remove items from cart
  - Clear cart
  
- **Orders**: Order processing and management
  - Checkout cart and create order
  - Get order history
  - Get order details
  - Get user certificates
  
- **Admin**: Administrative functions (requires authentication)
  - Get monthly reports
  - Get orders by date range
  - Get order statistics with growth metrics

### Regenerating Swagger Documentation

If you make changes to the API handlers, regenerate the Swagger docs:

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g main.go
```

## Data Persistence

- MongoDB data is persisted in the `mongodb_data` volume
- RabbitMQ data is persisted in the `rabbitmq_data` volume

## Troubleshooting

### Service Won't Start

1. Check if ports are already in use:
   ```bash
   lsof -i :8080
   lsof -i :27017
   lsof -i :5672
   ```

2. Check service logs:
   ```bash
   docker-compose logs order-service
   ```

### Database Connection Issues

1. Ensure MongoDB is running and accessible
2. Check the `MONGODB_URI` environment variable
3. Verify network connectivity between containers

### RabbitMQ Connection Issues

1. Ensure RabbitMQ is running and accessible
2. Check the `RABBITMQ_URL` environment variable
3. Verify the queue `certificate_generation` is created
