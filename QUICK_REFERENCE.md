# Carbon Clear - Quick Reference Guide

Quick commands and API endpoints for common operations.

## 🚀 Quick Start

```bash
# Start all services
./start.sh

# Or using Make
make build && make up

# Check status
make status

# View logs
make logs
```

## 📡 Service URLs

| Service | URL | Documentation |
|---------|-----|---------------|
| API Gateway | http://localhost:8000 | http://localhost:8000/swagger/index.html |
| User Service | http://localhost:8082 | http://localhost:8082/swagger/index.html |
| Project Service | http://localhost:8081 | http://localhost:8081/swagger/index.html |
| Order Service | http://localhost:8080 | http://localhost:8080/swagger/index.html |
| RabbitMQ | http://localhost:15672 | admin/admin123 |

## 🔐 Authentication

### Get User Token
```bash
TOKEN=$(curl -s -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.token')
```

### Get Admin Token
```bash
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8000/api/admin/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  | jq -r '.token')
```

### Use Token
```bash
curl http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

## 👤 User Operations

### Register User
```bash
curl -X POST http://localhost:8000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Profile
```bash
curl http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

## 🏭 Project Operations

### Browse Projects
```bash
curl http://localhost:8000/api/v1/projects
```

### Get Project Details
```bash
curl http://localhost:8000/api/v1/projects/1
```

### Search Projects
```bash
curl -X POST http://localhost:8000/api/v1/projects/search \
  -H "Content-Type: application/json" \
  -d '{
    "category": "Renewable Energy",
    "min_price": 10,
    "max_price": 50
  }'
```

### Get Categories/Regions/Countries
```bash
curl http://localhost:8000/api/v1/projects/categories
curl http://localhost:8000/api/v1/projects/regions
curl http://localhost:8000/api/v1/projects/countries
```

## 🛒 Cart Operations

### Add to Cart
```bash
curl -X POST http://localhost:8000/api/v1/cart/1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "1",
    "tons": 10
  }'
```

### View Cart
```bash
curl http://localhost:8000/api/v1/cart/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Update Cart Item
```bash
curl -X PUT http://localhost:8000/api/v1/cart/1/items/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"tons": 15}'
```

### Remove from Cart
```bash
curl -X DELETE http://localhost:8000/api/v1/cart/1/items/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Clear Cart
```bash
curl -X DELETE http://localhost:8000/api/v1/cart/1 \
  -H "Authorization: Bearer $TOKEN"
```

## 📦 Order Operations

### Checkout
```bash
curl -X POST http://localhost:8000/api/v1/orders/1/checkout \
  -H "Authorization: Bearer $TOKEN"
```

### Get Order History
```bash
curl http://localhost:8000/api/v1/orders/1/history \
  -H "Authorization: Bearer $TOKEN"
```

### Get Order Details
```bash
curl http://localhost:8000/api/v1/orders/ORDER_ID \
  -H "Authorization: Bearer $TOKEN"
```

### Get Certificates
```bash
curl http://localhost:8000/api/v1/orders/1/certificates \
  -H "Authorization: Bearer $TOKEN"
```

## 👨‍💼 Admin Operations

### Create Project
```bash
curl -X POST http://localhost:8000/api/v1/projects/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Solar Farm",
    "description": "Solar energy project",
    "category": "Renewable Energy",
    "location": "California",
    "region": "North America",
    "country": "USA",
    "price_per_ton": 25.00,
    "available_tons": 1000,
    "status": "active"
  }'
```

### Update Project
```bash
curl -X PUT http://localhost:8000/api/v1/projects/admin/1 \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "price_per_ton": 30.00,
    "available_tons": 800
  }'
```

### Delete Project
```bash
curl -X DELETE http://localhost:8000/api/v1/projects/admin/1 \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Get All Users
```bash
curl http://localhost:8000/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Get Monthly Report
```bash
curl "http://localhost:8000/api/v1/admin/reports/monthly?year=2024&month=1" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Get Statistics
```bash
curl http://localhost:8000/api/v1/admin/statistics \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 🐳 Docker Commands

### Start Services
```bash
docker-compose up -d              # Start all services
docker-compose up -d api_gateway  # Start specific service
```

### Stop Services
```bash
docker-compose down              # Stop all services
docker-compose stop api_gateway  # Stop specific service
```

### View Logs
```bash
docker-compose logs -f                    # All services
docker-compose logs -f api_gateway        # Specific service
docker-compose logs --tail=100 api_gateway # Last 100 lines
```

### Restart Services
```bash
docker-compose restart                # All services
docker-compose restart api_gateway    # Specific service
```

### Rebuild Services
```bash
docker-compose up -d --build          # Rebuild and start
docker-compose build api_gateway      # Rebuild specific service
```

### Check Status
```bash
docker-compose ps                     # Service status
docker stats                          # Resource usage
```

### Clean Up
```bash
docker-compose down -v                # Stop and remove volumes
docker-compose down --rmi all         # Stop and remove images
docker system prune -a                # Clean all unused resources
```

## 🗄️ Database Access

### PostgreSQL (User Service)
```bash
docker exec -it postgres_user psql -U postgres -d user_service_db

# Common queries
SELECT * FROM users;
\dt                    # List tables
\d users              # Describe users table
```

### PostgreSQL (Project Service)
```bash
docker exec -it postgres_project psql -U postgres -d project_service_db

# Common queries
SELECT * FROM projects;
SELECT * FROM projects WHERE category = 'Renewable Energy';
```

### MongoDB (Order Service)
```bash
docker exec -it mongodb mongosh -u admin -p admin123 --authenticationDatabase admin

# Common commands
use order_service_db
show collections
db.orders.find().pretty()
db.carts.find().pretty()
db.certificates.find().pretty()
```

### Redis
```bash
docker exec -it redis redis-cli -a redis123

# Common commands
KEYS *
GET project:list:*
FLUSHALL              # Clear all cache (careful!)
```

## 📊 Monitoring

### Health Checks
```bash
curl http://localhost:8000/health     # API Gateway
curl http://localhost:8082/           # User Service
curl http://localhost:8081/health     # Project Service
curl http://localhost:8080/swagger/   # Order Service
```

### Service Metrics
```bash
# Request count
docker-compose logs api_gateway | grep -c "GET"

# Error count
docker-compose logs api_gateway | grep -c "error"

# Response times
docker-compose logs api_gateway | grep "Latency"
```

### Resource Usage
```bash
docker stats                          # Real-time stats
docker stats --no-stream              # One-time snapshot
```

## 🛠️ Make Commands

```bash
make help            # Show all commands
make build           # Build all images
make up              # Start all services
make down            # Stop all services
make logs            # View all logs
make restart         # Restart all services
make status          # Check service status
make test            # Run health checks
make clean           # Remove everything

# Individual services
make gateway-logs
make gateway-restart
make user-logs
make user-restart
make project-logs
make project-restart
make order-logs
make order-restart
```

## 🧪 Testing

### Import Postman Collection
1. Open Postman
2. Import `carbon-clear-api-collection.json`
3. Collection includes all endpoints with examples

### Run Test Suite
```bash
chmod +x test-suite.sh
./test-suite.sh
```

### Load Testing
```bash
# Install apache bench
brew install apache-bench    # macOS
apt-get install apache2-utils # Ubuntu

# Run load test
ab -n 1000 -c 10 http://localhost:8000/health
```

## 🔍 Troubleshooting

### Services Won't Start
```bash
# Check logs
docker-compose logs

# Check ports
lsof -i :8000
lsof -i :8080
lsof -i :8081
lsof -i :8082

# Restart from scratch
make down
make clean
make build
make up
```

### Authentication Errors
```bash
# Verify token
echo $TOKEN

# Get fresh token
TOKEN=$(curl -s -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.token')
```

### Database Connection Errors
```bash
# Check database health
docker-compose ps | grep postgres
docker-compose ps | grep mongodb

# Restart databases
docker-compose restart postgres_user
docker-compose restart postgres_project
docker-compose restart mongodb
```

### Rate Limit Errors
```bash
# Wait 1 minute or adjust rate limit
# Edit .env
RATE_LIMIT_PER_MIN=200

# Restart gateway
docker-compose restart api_gateway
```

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [README.md](README.md) | Overview and quick start |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System architecture details |
| [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) | Deployment instructions |
| [TESTING_GUIDE.md](TESTING_GUIDE.md) | Testing procedures |
| [api_gateway/README.md](api_gateway/README.md) | API Gateway documentation |

## 🔗 Useful Links

- **Swagger UI**: http://localhost:8000/swagger/index.html
- **RabbitMQ Management**: http://localhost:15672
- **Health Status**: http://localhost:8000/health

## 💡 Tips & Tricks

### Save Frequently Used Tokens
```bash
# Add to ~/.bashrc or ~/.zshrc
export CARBON_TOKEN="your-token-here"
export CARBON_ADMIN_TOKEN="your-admin-token-here"

# Use in commands
curl -H "Authorization: Bearer $CARBON_TOKEN" ...
```

### Create Aliases
```bash
# Add to ~/.bashrc or ~/.zshrc
alias cc-up='cd ~/carbon-clear && make up'
alias cc-down='cd ~/carbon-clear && make down'
alias cc-logs='cd ~/carbon-clear && make logs'
alias cc-status='cd ~/carbon-clear && make status'
```

### Quick Test Script
```bash
#!/bin/bash
# Save as test.sh

# Register and login
RESPONSE=$(curl -s -X POST http://localhost:8000/api/users/register \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test$(date +%s)\",\"email\":\"test$(date +%s)@example.com\",\"password\":\"password\"}")

TOKEN=$(echo $RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

# Test authenticated endpoint
curl -s http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

## 📞 Support

- **Email**: support@carbonclear.com
- **Issues**: GitHub Issues
- **Documentation**: See docs folder

---

**Last Updated**: 2024-01-01
**Version**: 1.0.0

