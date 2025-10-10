# API Gateway - Complete Route Map

This document provides a complete visual map of all routes handled by the API Gateway.

## Route Flow Diagram

```
                              API Gateway (Port 8000)
                                        |
        ┌───────────────────────────────┼───────────────────────────────┐
        |                               |                               |
        |                               |                               |
   User Routes                    Project Routes                  Order Routes
   (→ :8082)                       (→ :8081)                       (→ :8080)
```

## Public Routes (No Authentication Required)

### User Service Routes

```
POST /api/users/register
├─ Description: Register new user
├─ Body: { name, email, password }
└─ Response: { message, token, user }

POST /api/users/login
├─ Description: User login
├─ Body: { email, password }
└─ Response: { token, user }

POST /api/admin/users/register
├─ Description: Register admin user
├─ Body: { name, email, password }
└─ Response: { message, token, user }

POST /api/admin/users/login
├─ Description: Admin login
├─ Body: { email, password }
└─ Response: { token, user }
```

### Project Service Routes

```
GET /api/v1/projects
├─ Description: List all projects
├─ Query Params: page, limit, status
└─ Response: { projects[], total, page, limit }

GET /api/v1/projects/:id
├─ Description: Get project details
├─ Params: id (project ID)
└─ Response: { project }

POST /api/v1/projects/search
├─ Description: Search/filter projects
├─ Body: { category, region, country, min_price, max_price }
└─ Response: { projects[], total }

GET /api/v1/projects/categories
├─ Description: Get available categories
└─ Response: { categories[] }

GET /api/v1/projects/regions
├─ Description: Get available regions
└─ Response: { regions[] }

GET /api/v1/projects/countries
├─ Description: Get available countries
└─ Response: { countries[] }
```

## Protected Routes (User JWT Required)

### User Profile

```
GET /api/users/profile
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Get user profile
└─ Response: { user }
```

### Shopping Cart

```
POST /api/v1/cart/:userID/items
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Add item to cart
├─ Params: userID
├─ Body: { project_id, tons }
└─ Response: { message, cart }

GET /api/v1/cart/:userID
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Get user's cart
├─ Params: userID
└─ Response: { cart }

PUT /api/v1/cart/:userID/items/:projectID
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Update cart item quantity
├─ Params: userID, projectID
├─ Body: { tons }
└─ Response: { message, cart }

DELETE /api/v1/cart/:userID/items/:projectID
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Remove item from cart
├─ Params: userID, projectID
└─ Response: { message }

DELETE /api/v1/cart/:userID
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Clear entire cart
├─ Params: userID
└─ Response: { message }
```

### Orders

```
POST /api/v1/orders/:userID/checkout
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Checkout cart and create order
├─ Params: userID
└─ Response: { message, order_id, total_amount }

GET /api/v1/orders/:userID/history
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Get user's order history
├─ Params: userID
├─ Query: page, limit
└─ Response: { orders[], total }

GET /api/v1/orders/:orderID
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Get order details
├─ Params: orderID
└─ Response: { order }

GET /api/v1/orders/:userID/certificates
├─ Auth: Bearer {USER_TOKEN}
├─ Description: Get user's certificates
├─ Params: userID
└─ Response: { certificates[] }
```

## Protected Routes (Admin JWT Required)

### Admin User Management

```
GET /api/admin/users
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Get all users
├─ Query: page, limit, role
└─ Response: { users[], total }

POST /api/admin/users
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Create new user
├─ Body: { name, email, password, role }
└─ Response: { message, user }

GET /api/admin/users/:id
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Get user by ID
├─ Params: id
└─ Response: { user }

PUT /api/admin/users/:id
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Update user
├─ Params: id
├─ Body: { name, email, role }
└─ Response: { message, user }

DELETE /api/admin/users/:id
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Delete user
├─ Params: id
└─ Response: { message }
```

### Admin Project Management

```
POST /api/v1/projects/admin
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Create new project
├─ Body: { name, description, category, location, region, 
│          country, price_per_ton, available_tons, 
│          image_url, status }
└─ Response: { message, project }

PUT /api/v1/projects/admin/:id
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Update project
├─ Params: id
├─ Body: { any project fields to update }
└─ Response: { message, project }

DELETE /api/v1/projects/admin/:id
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Delete project
├─ Params: id
└─ Response: { message }
```

### Admin Reports & Statistics

```
GET /api/v1/admin/reports/monthly
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Get monthly sales report
├─ Query: year, month
└─ Response: { report }

GET /api/v1/admin/orders/date-range
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Get orders in date range
├─ Query: start, end (YYYY-MM-DD)
└─ Response: { orders[], total, revenue }

GET /api/v1/admin/statistics
├─ Auth: Bearer {ADMIN_TOKEN}
├─ Description: Get overall statistics
└─ Response: { 
    total_orders, 
    total_revenue, 
    total_tons_sold,
    total_users 
  }
```

## Gateway Special Routes

```
GET /health
├─ Auth: None
├─ Description: Gateway health check
└─ Response: { status: "healthy", service: "api_gateway", version: "1.0.0" }

GET /
├─ Auth: None
├─ Description: Gateway information
└─ Response: { 
    message: "Carbon Clear API Gateway",
    version: "1.0.0",
    services: { ... }
  }

GET /swagger/*
├─ Auth: None
├─ Description: Interactive API documentation
└─ Response: Swagger UI
```

## Authentication Flow

```
┌──────────────┐
│   Client     │
└──────┬───────┘
       │ 1. POST /api/users/login
       │    { email, password }
       ▼
┌──────────────┐
│  API Gateway │
└──────┬───────┘
       │ 2. Forward to User Service
       ▼
┌──────────────┐
│ User Service │ 3. Validate credentials
└──────┬───────┘ 4. Generate JWT token
       │
       │ 5. Return token
       ▼
┌──────────────┐
│  API Gateway │ 6. Forward response
└──────┬───────┘
       │
       │ 7. Return: { token, user }
       ▼
┌──────────────┐
│   Client     │ 8. Store token
└──────┬───────┘
       │
       │ 9. Use token in future requests
       │    Header: Authorization: Bearer {token}
       ▼
┌──────────────┐
│  API Gateway │ 10. Validate JWT
└──────┬───────┘ 11. If valid, forward request
       │
       ▼
```

## Request Middleware Chain

```
Incoming Request
       │
       ▼
┌─────────────────┐
│ Custom Logger   │ Log: method, path, IP
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Rate Limiter    │ Check: requests/minute per IP
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ CORS Handler    │ Set: CORS headers
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ JWT Middleware  │ Validate: token if protected route
│ (if protected)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Proxy Handler   │ Forward: to backend service
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Backend Service │ Process: request
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Proxy Handler   │ Return: response to client
└────────┬────────┘
         │
         ▼
Response to Client
```

## Error Responses

All routes can return these error responses:

```
400 Bad Request
├─ Invalid input format
└─ Missing required fields

401 Unauthorized
├─ Missing authentication token
├─ Invalid token
└─ Expired token

403 Forbidden
├─ Insufficient permissions
└─ Wrong token type (user vs admin)

404 Not Found
├─ Resource doesn't exist
└─ Invalid endpoint

429 Too Many Requests
├─ Rate limit exceeded
└─ Try again after 1 minute

500 Internal Server Error
├─ Server error
└─ Database error

502 Bad Gateway
├─ Backend service unavailable
└─ Backend service timeout

503 Service Unavailable
├─ Gateway overloaded
└─ Maintenance mode
```

## Rate Limiting

```
Rate Limit Configuration:
├─ Default: 100 requests/minute per IP
├─ Burst: 10 requests (10% of rate)
├─ Cleanup: Every 5 minutes
└─ Response: 429 when exceeded

Headers in Response:
├─ X-RateLimit-Limit: 100
├─ X-RateLimit-Remaining: 95
└─ X-RateLimit-Reset: 1704067200
```

## CORS Configuration

```
Allowed Origins: * (configurable)
Allowed Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
Allowed Headers: Origin, Content-Type, Accept, Authorization
Exposed Headers: Content-Length
Max Age: 3600
```

## Route Priority

Routes are matched in this order:

1. **Exact matches** (e.g., `/health`)
2. **Static paths** (e.g., `/api/users/register`)
3. **Path parameters** (e.g., `/api/users/:id`)
4. **Wildcard** (e.g., `/swagger/*`)

## Service Discovery Map

```yaml
User Service (http://localhost:8082):
  - /api/users/*
  - /api/admin/users/*

Project Service (http://localhost:8081):
  - /api/v1/projects/*

Order Service (http://localhost:8080):
  - /api/v1/cart/*
  - /api/v1/orders/*
  - /api/v1/admin/reports/*
  - /api/v1/admin/orders/*
  - /api/v1/admin/statistics
```

## Complete Endpoint Count

| Category | Count | Auth Type |
|----------|-------|-----------|
| Health & Info | 3 | None |
| User Auth | 4 | None |
| User Profile | 1 | User JWT |
| Admin Users | 5 | Admin JWT |
| Projects (Public) | 6 | None |
| Projects (Admin) | 3 | Admin JWT |
| Shopping Cart | 5 | User JWT |
| Orders | 4 | User JWT |
| Admin Reports | 3 | Admin JWT |
| **TOTAL** | **34** | - |

## Testing Workflow

```bash
# 1. Health check (no auth)
curl http://localhost:8000/health

# 2. Register user (no auth)
curl -X POST http://localhost:8000/api/users/register \
  -d '{"name":"User","email":"user@test.com","password":"pass"}'

# 3. Login user (no auth, get token)
TOKEN=$(curl -X POST http://localhost:8000/api/users/login \
  -d '{"email":"user@test.com","password":"pass"}' | jq -r '.token')

# 4. Access protected route (with token)
curl http://localhost:8000/api/users/profile \
  -H "Authorization: Bearer $TOKEN"

# 5. Browse projects (no auth)
curl http://localhost:8000/api/v1/projects

# 6. Add to cart (with token)
curl -X POST http://localhost:8000/api/v1/cart/1/items \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"project_id":"1","tons":10}'

# 7. Checkout (with token)
curl -X POST http://localhost:8000/api/v1/orders/1/checkout \
  -H "Authorization: Bearer $TOKEN"

# 8. Get certificates (with token)
curl http://localhost:8000/api/v1/orders/1/certificates \
  -H "Authorization: Bearer $TOKEN"
```

---

**Note**: Replace `:userID`, `:projectID`, `:orderID`, and `:id` with actual IDs when making requests.

**Tip**: Use the Postman collection (`carbon-clear-api-collection.json`) for easier testing with automatic token management.

