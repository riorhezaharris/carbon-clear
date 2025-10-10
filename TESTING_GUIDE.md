# Carbon Clear - API Testing Guide

This guide provides comprehensive instructions for testing the Carbon Clear API Gateway and microservices.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Using Postman/Thunder Client](#using-postmanthunder-client)
3. [Manual Testing with cURL](#manual-testing-with-curl)
4. [Testing Workflows](#testing-workflows)
5. [Common Test Scenarios](#common-test-scenarios)
6. [Troubleshooting](#troubleshooting)

## Getting Started

### Prerequisites

1. **Start all services:**
   ```bash
   ./start.sh
   # or
   make up
   ```

2. **Verify services are running:**
   ```bash
   make test
   ```

### API Base URL

All requests go through the API Gateway:
```
http://localhost:8000
```

## Using Postman/Thunder Client

### Import the Collection

1. Open Postman or Thunder Client
2. Import the collection file: `carbon-clear-api-collection.json`
3. The collection includes all endpoints with example requests

### Collection Variables

The collection uses variables for easy testing:

- `baseUrl`: http://localhost:8000
- `userToken`: Automatically set after login
- `adminToken`: Automatically set after admin login
- `userId`: User ID for requests
- `projectId`: Project ID for cart/orders
- `orderId`: Order ID after checkout

### Running the Collection

1. **Register a new user** → automatically saves `userToken`
2. **Login** → updates `userToken`
3. Use other endpoints with the saved token

## Manual Testing with cURL

### Health Checks

```bash
# API Gateway health
curl http://localhost:8000/health

# Gateway info
curl http://localhost:8000/
```

### User Registration and Login

#### Register a User

```bash
curl -X POST http://localhost:8000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

Expected Response:
```json
{
  "message": "User registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

#### Login User

```bash
curl -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Save the token:
```bash
TOKEN=$(curl -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}' \
  | jq -r '.token')
```

#### Get Profile

```bash
curl http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Admin Operations

#### Register Admin

```bash
curl -X POST http://localhost:8000/api/admin/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@carbonclear.com",
    "password": "admin123"
  }'
```

#### Login Admin

```bash
ADMIN_TOKEN=$(curl -X POST http://localhost:8000/api/admin/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@carbonclear.com","password":"admin123"}' \
  | jq -r '.token')
```

#### Get All Users

```bash
curl http://localhost:8000/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Project Management

#### Browse Projects (No Auth Required)

```bash
curl http://localhost:8000/api/v1/projects
```

#### Get Project Details

```bash
curl http://localhost:8000/api/v1/projects/1
```

#### Search Projects

```bash
curl -X POST http://localhost:8000/api/v1/projects/search \
  -H "Content-Type: application/json" \
  -d '{
    "category": "Renewable Energy",
    "region": "Asia",
    "min_price": 10,
    "max_price": 50
  }'
```

#### Get Categories, Regions, Countries

```bash
curl http://localhost:8000/api/v1/projects/categories
curl http://localhost:8000/api/v1/projects/regions
curl http://localhost:8000/api/v1/projects/countries
```

#### Create Project (Admin Only)

```bash
curl -X POST http://localhost:8000/api/v1/projects/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Solar Farm Initiative",
    "description": "Large-scale solar energy project",
    "category": "Renewable Energy",
    "location": "California, USA",
    "region": "North America",
    "country": "USA",
    "price_per_ton": 25.00,
    "available_tons": 1000,
    "image_url": "https://example.com/solar.jpg",
    "status": "active"
  }'
```

### Shopping Cart

#### Add to Cart

```bash
curl -X POST http://localhost:8000/api/v1/cart/1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "1",
    "tons": 10
  }'
```

#### View Cart

```bash
curl http://localhost:8000/api/v1/cart/1 \
  -H "Authorization: Bearer $TOKEN"
```

#### Update Cart Item

```bash
curl -X PUT http://localhost:8000/api/v1/cart/1/items/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tons": 15
  }'
```

#### Remove Item from Cart

```bash
curl -X DELETE http://localhost:8000/api/v1/cart/1/items/1 \
  -H "Authorization: Bearer $TOKEN"
```

#### Clear Cart

```bash
curl -X DELETE http://localhost:8000/api/v1/cart/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Orders

#### Checkout

```bash
curl -X POST http://localhost:8000/api/v1/orders/1/checkout \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Order History

```bash
curl http://localhost:8000/api/v1/orders/1/history \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Order Details

```bash
curl http://localhost:8000/api/v1/orders/ORDER_ID \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Certificates

```bash
curl http://localhost:8000/api/v1/orders/1/certificates \
  -H "Authorization: Bearer $TOKEN"
```

### Admin Reports

#### Get Monthly Report

```bash
curl "http://localhost:8000/api/v1/admin/reports/monthly?year=2024&month=1" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### Get Orders by Date Range

```bash
curl "http://localhost:8000/api/v1/admin/orders/date-range?start=2024-01-01&end=2024-12-31" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### Get Statistics

```bash
curl http://localhost:8000/api/v1/admin/statistics \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## Testing Workflows

### Complete User Flow

```bash
#!/bin/bash

# 1. Register user
echo "1. Registering user..."
RESPONSE=$(curl -s -X POST http://localhost:8000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }')

TOKEN=$(echo $RESPONSE | jq -r '.token')
USER_ID=$(echo $RESPONSE | jq -r '.user.id')
echo "Token: $TOKEN"
echo "User ID: $USER_ID"

# 2. Browse projects
echo -e "\n2. Browsing projects..."
curl -s http://localhost:8000/api/v1/projects | jq '.'

# 3. Add to cart
echo -e "\n3. Adding to cart..."
curl -s -X POST http://localhost:8000/api/v1/cart/$USER_ID/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "1",
    "tons": 10
  }' | jq '.'

# 4. View cart
echo -e "\n4. Viewing cart..."
curl -s http://localhost:8000/api/v1/cart/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 5. Checkout
echo -e "\n5. Checking out..."
ORDER=$(curl -s -X POST http://localhost:8000/api/v1/orders/$USER_ID/checkout \
  -H "Authorization: Bearer $TOKEN")
echo $ORDER | jq '.'

ORDER_ID=$(echo $ORDER | jq -r '.order_id')

# 6. Get order history
echo -e "\n6. Getting order history..."
curl -s http://localhost:8000/api/v1/orders/$USER_ID/history \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 7. View certificates (after a few seconds for generation)
echo -e "\n7. Waiting for certificate generation..."
sleep 5
curl -s http://localhost:8000/api/v1/orders/$USER_ID/certificates \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

### Admin Workflow

```bash
#!/bin/bash

# 1. Register admin
echo "1. Registering admin..."
ADMIN_RESPONSE=$(curl -s -X POST http://localhost:8000/api/admin/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@carbonclear.com",
    "password": "admin123"
  }')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.token')
echo "Admin Token: $ADMIN_TOKEN"

# 2. Create project
echo -e "\n2. Creating project..."
curl -s -X POST http://localhost:8000/api/v1/projects/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wind Energy Project",
    "description": "Offshore wind farm",
    "category": "Renewable Energy",
    "location": "North Sea",
    "region": "Europe",
    "country": "Netherlands",
    "price_per_ton": 30.00,
    "available_tons": 500,
    "image_url": "https://example.com/wind.jpg",
    "status": "active"
  }' | jq '.'

# 3. Get all users
echo -e "\n3. Getting all users..."
curl -s http://localhost:8000/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'

# 4. Get statistics
echo -e "\n4. Getting statistics..."
curl -s http://localhost:8000/api/v1/admin/statistics \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'

# 5. Get monthly report
echo -e "\n5. Getting monthly report..."
curl -s "http://localhost:8000/api/v1/admin/reports/monthly?year=2024&month=1" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
```

## Common Test Scenarios

### Test Authentication

```bash
# Should return 401 without token
curl http://localhost:8000/api/users/profile

# Should return 401 with invalid token
curl http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer invalid-token"

# Should succeed with valid token
curl http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Test Rate Limiting

```bash
# Send 150 requests quickly (rate limit is 100/min)
for i in {1..150}; do
  curl -s http://localhost:8000/health > /dev/null
  echo "Request $i"
done

# Should start getting 429 Too Many Requests after 100 requests
```

### Test CORS

```bash
# Test preflight request
curl -X OPTIONS http://localhost:8000/api/v1/projects \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

### Test Error Handling

```bash
# Invalid JSON
curl -X POST http://localhost:8000/api/users/register \
  -H "Content-Type: application/json" \
  -d 'invalid json'

# Missing required fields
curl -X POST http://localhost:8000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{}'

# Non-existent resource
curl http://localhost:8000/api/v1/projects/99999
```

## Troubleshooting

### Common Issues

#### 401 Unauthorized

**Problem:** Getting 401 errors on protected routes

**Solutions:**
1. Verify token is included in Authorization header
2. Check token format: `Bearer <token>`
3. Ensure token hasn't expired - login again
4. Verify correct token type (user vs admin)

```bash
# Check token
echo $TOKEN

# Re-login
TOKEN=$(curl -X POST http://localhost:8000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.token')
```

#### 429 Too Many Requests

**Problem:** Rate limit exceeded

**Solutions:**
1. Wait 1 minute before retrying
2. Adjust rate limit in gateway configuration
3. Use different IP for testing

```bash
# Check rate limit configuration
docker-compose exec api_gateway env | grep RATE_LIMIT
```

#### 502 Bad Gateway

**Problem:** Gateway can't connect to backend service

**Solutions:**
1. Check if all services are running
2. Verify service health

```bash
# Check service status
docker-compose ps

# Check logs
docker-compose logs user_service
docker-compose logs project_service
docker-compose logs order_service

# Restart services
docker-compose restart
```

#### Connection Refused

**Problem:** Can't connect to API Gateway

**Solutions:**
1. Verify gateway is running
2. Check port 8000 is not in use

```bash
# Check if gateway is running
docker-compose ps api_gateway

# Check port
lsof -i :8000

# Restart gateway
docker-compose restart api_gateway
```

### Debugging Tips

#### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api_gateway

# Last 100 lines
docker-compose logs --tail=100 api_gateway
```

#### Check Service Health

```bash
# Health checks
curl http://localhost:8000/health
curl http://localhost:8082/
curl http://localhost:8081/health

# Detailed status
docker-compose ps
```

#### Test Direct Service Access

```bash
# Bypass gateway - test service directly
curl http://localhost:8082/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

## Performance Testing

### Using Apache Bench

```bash
# Install ab
# macOS: brew install apache-bench
# Ubuntu: apt-get install apache2-utils

# Test health endpoint
ab -n 1000 -c 10 http://localhost:8000/health

# Test with authentication
ab -n 1000 -c 10 -H "Authorization: Bearer $TOKEN" \
  http://localhost:8000/api/users/profile
```

### Using wrk

```bash
# Install wrk
# macOS: brew install wrk
# Ubuntu: apt-get install wrk

# Basic test
wrk -t12 -c400 -d30s http://localhost:8000/health

# With custom header
wrk -t12 -c400 -d30s \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8000/api/users/profile
```

## Automated Testing

### Create Test Suite

Save as `test-suite.sh`:

```bash
#!/bin/bash

PASSED=0
FAILED=0

test_health() {
  echo "Testing health endpoint..."
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8000/health)
  if [ "$STATUS" = "200" ]; then
    echo "✓ Health check passed"
    ((PASSED++))
  else
    echo "✗ Health check failed (Status: $STATUS)"
    ((FAILED++))
  fi
}

test_register() {
  echo "Testing user registration..."
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST http://localhost:8000/api/users/register \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Test$(date +%s)\",\"email\":\"test$(date +%s)@example.com\",\"password\":\"password123\"}")
  
  if [ "$STATUS" = "201" ]; then
    echo "✓ Registration passed"
    ((PASSED++))
  else
    echo "✗ Registration failed (Status: $STATUS)"
    ((FAILED++))
  fi
}

# Run tests
test_health
test_register

# Summary
echo ""
echo "Test Summary:"
echo "Passed: $PASSED"
echo "Failed: $FAILED"
```

Run:
```bash
chmod +x test-suite.sh
./test-suite.sh
```

## Resources

- **Swagger Documentation**: http://localhost:8000/swagger/index.html
- **Collection File**: `carbon-clear-api-collection.json`
- **Deployment Guide**: `DEPLOYMENT_GUIDE.md`
- **Main README**: `README.md`

## Support

If you encounter issues:
1. Check the logs: `docker-compose logs`
2. Verify configuration: `.env` file
3. Review documentation: README files
4. Contact: support@carbonclear.com

