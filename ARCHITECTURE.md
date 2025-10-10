# Carbon Clear - System Architecture

This document provides a comprehensive overview of the Carbon Clear microservices architecture.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [System Components](#system-components)
3. [API Gateway](#api-gateway)
4. [Microservices](#microservices)
5. [Data Layer](#data-layer)
6. [Communication Patterns](#communication-patterns)
7. [Security Architecture](#security-architecture)
8. [Scalability & Performance](#scalability--performance)

## Architecture Overview

Carbon Clear follows a **microservices architecture** pattern with an API Gateway as the single entry point for all client communications.

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│          (Web App, Mobile App, Third-party Services)        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway (Port 8000)                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  • Authentication & Authorization                      │  │
│  │  • Rate Limiting (100 req/min)                        │  │
│  │  • Request Routing                                    │  │
│  │  • CORS Handling                                      │  │
│  │  • Logging & Monitoring                               │  │
│  └──────────────────────────────────────────────────────┘  │
└─────┬─────────────────┬─────────────────┬──────────────────┘
      │                 │                 │
      ▼                 ▼                 ▼
┌──────────┐      ┌──────────┐     ┌──────────┐
│   User   │      │ Project  │     │  Order   │
│ Service  │      │ Service  │     │ Service  │
│ :8082    │      │ :8081    │     │ :8080    │
└────┬─────┘      └────┬─────┘     └────┬─────┘
     │                 │                 │
     ▼                 ▼                 ▼
┌──────────┐      ┌──────────┐     ┌──────────┐
│PostgreSQL│      │PostgreSQL│     │ MongoDB  │
│  :5433   │      │  :5434   │     │  :27017  │
└──────────┘      └────┬─────┘     └────┬─────┘
                       │                 │
                  ┌────┴────┐       ┌────┴────┐
                  ▼         ▼       ▼         ▼
              ┌─────┐   ┌─────┐ ┌────────┐ ┌─────┐
              │Redis│   │ ES  │ │RabbitMQ│ │Cache│
              │:6379│   │:9200│ │ :5672  │ └─────┘
              └─────┘   └─────┘ └────────┘
```

## System Components

### 1. API Gateway
- **Technology**: Go + Echo Framework
- **Port**: 8000
- **Responsibilities**:
  - Single entry point for all client requests
  - JWT authentication and authorization
  - Rate limiting and throttling
  - Request routing to backend services
  - CORS policy enforcement
  - Centralized logging

### 2. User Service
- **Technology**: Go + Echo + PostgreSQL
- **Port**: 8082
- **Database**: PostgreSQL (port 5433)
- **Responsibilities**:
  - User registration and authentication
  - Admin management
  - Profile management
  - JWT token generation and validation

### 3. Project Service
- **Technology**: Go + Echo + PostgreSQL + Redis + Elasticsearch
- **Port**: 8081
- **Database**: PostgreSQL (port 5434)
- **Cache**: Redis (port 6379)
- **Search**: Elasticsearch (port 9200)
- **Responsibilities**:
  - Carbon offset project catalog
  - Project search and filtering
  - Category and region management
  - Admin project CRUD operations
  - Redis caching for performance
  - Full-text search with Elasticsearch

### 4. Order Service
- **Technology**: Go + Echo + MongoDB + RabbitMQ
- **Port**: 8080
- **Database**: MongoDB (port 27017)
- **Message Queue**: RabbitMQ (port 5672)
- **Responsibilities**:
  - Shopping cart management
  - Order processing and checkout
  - Certificate generation (async via RabbitMQ)
  - Order history and tracking
  - Monthly reports and statistics
  - Scheduled tasks

## API Gateway

### Design Patterns

#### 1. Gateway Routing Pattern
Routes requests to appropriate backend services based on URL patterns:

```go
/api/users/*           → User Service
/api/v1/projects/*     → Project Service
/api/v1/cart/*         → Order Service
/api/v1/orders/*       → Order Service
/api/v1/admin/*        → Order Service (admin endpoints)
```

#### 2. Gateway Authentication Pattern
Centralizes authentication logic:

```go
// User JWT validation
UserJWTMiddleware → validates token → forwards to service

// Admin JWT validation
AdminJWTMiddleware → validates token → forwards to service
```

#### 3. Rate Limiting Pattern
Implements token bucket algorithm:
- **Default**: 100 requests per minute per IP
- **Burst**: 10% of rate limit
- **Cleanup**: Every 5 minutes

### Request Flow

```
Client Request
    ↓
[Custom Logger] - Log request details
    ↓
[Rate Limiter] - Check rate limit
    ↓
[CORS Middleware] - Handle CORS
    ↓
[Auth Middleware] - Validate JWT (if protected route)
    ↓
[Proxy Handler] - Forward to backend service
    ↓
[Response] - Return to client
```

## Microservices

### User Service Architecture

```
┌────────────────────────────────────┐
│         User Service               │
│                                    │
│  ┌──────────────────────────┐    │
│  │    HTTP Handlers          │    │
│  │  • RegisterUser           │    │
│  │  • LoginUser              │    │
│  │  • GetProfile             │    │
│  │  • UpdateUser             │    │
│  └───────────┬──────────────┘    │
│              ▼                     │
│  ┌──────────────────────────┐    │
│  │    Repository Layer       │    │
│  │  • UserRepository         │    │
│  └───────────┬──────────────┘    │
│              ▼                     │
│  ┌──────────────────────────┐    │
│  │    PostgreSQL Database    │    │
│  │  Table: users             │    │
│  └──────────────────────────┘    │
└────────────────────────────────────┘
```

**Database Schema:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Project Service Architecture

```
┌─────────────────────────────────────────────┐
│          Project Service                     │
│                                              │
│  ┌──────────────────────────────┐           │
│  │    HTTP Handlers              │           │
│  │  • GetAllProjects             │           │
│  │  • SearchProjects             │           │
│  │  • CreateProject (admin)      │           │
│  └───────────┬──────────────────┘           │
│              ▼                                │
│  ┌──────────────────────────────┐           │
│  │    Repository Layer           │           │
│  │  • ProjectRepository          │           │
│  │  • ElasticsearchClient        │           │
│  └───────────┬──────────────────┘           │
│              ▼                                │
│  ┌──────────────────┬───────────────┐       │
│  │   PostgreSQL     │  Elasticsearch│       │
│  │   (Primary DB)   │  (Search)     │       │
│  └──────────────────┴───────────────┘       │
│              ▼                                │
│  ┌──────────────────────────────┐           │
│  │    Redis Cache                │           │
│  │  • Cache project lists        │           │
│  │  • TTL: configurable          │           │
│  └──────────────────────────────┘           │
└─────────────────────────────────────────────┘
```

**Database Schema:**
```sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    location VARCHAR(255),
    region VARCHAR(100),
    country VARCHAR(100),
    price_per_ton DECIMAL(10,2),
    available_tons INTEGER,
    image_url TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Order Service Architecture

```
┌──────────────────────────────────────────────────┐
│              Order Service                        │
│                                                   │
│  ┌──────────────────────────────┐                │
│  │    HTTP Handlers              │                │
│  │  • AddToCart                  │                │
│  │  • Checkout                   │                │
│  │  • GetOrderHistory            │                │
│  └───────────┬──────────────────┘                │
│              ▼                                     │
│  ┌──────────────────────────────┐                │
│  │    Repository Layer           │                │
│  │  • CartRepository             │                │
│  │  • OrderRepository            │                │
│  │  • CertificateRepository      │                │
│  └───────────┬──────────────────┘                │
│              ▼                                     │
│  ┌──────────────────────────────┐                │
│  │    MongoDB Database           │                │
│  │  • carts collection           │                │
│  │  • orders collection          │                │
│  │  • certificates collection    │                │
│  └──────────────────────────────┘                │
│              ▼                                     │
│  ┌──────────────────────────────┐                │
│  │    Services Layer             │                │
│  │  • CertificateService         │                │
│  │  • SchedulerService           │                │
│  └───────────┬──────────────────┘                │
│              ▼                                     │
│  ┌──────────────────────────────┐                │
│  │    RabbitMQ                   │                │
│  │  Queue: certificate_queue     │                │
│  │  • Publish: on order creation │                │
│  │  • Consume: generate cert     │                │
│  └──────────────────────────────┘                │
└──────────────────────────────────────────────────┘
```

**MongoDB Collections:**
```javascript
// Carts
{
  _id: ObjectId,
  user_id: "1",
  items: [
    {
      project_id: "1",
      tons: 10,
      price_per_ton: 25.00
    }
  ],
  created_at: ISODate,
  updated_at: ISODate
}

// Orders
{
  _id: ObjectId,
  order_id: "ORD-20240101-ABC123",
  user_id: "1",
  items: [...],
  total_amount: 250.00,
  status: "completed",
  created_at: ISODate
}

// Certificates
{
  _id: ObjectId,
  certificate_number: "CERT-20240101-XYZ789",
  order_id: "ORD-20240101-ABC123",
  user_id: "1",
  tons_offset: 10,
  project_details: {...},
  generated_at: ISODate
}
```

## Data Layer

### Database Selection Rationale

#### PostgreSQL (User & Project Services)
**Why PostgreSQL?**
- ACID compliance for user data integrity
- Strong support for relational data (users, projects)
- Excellent query performance with indexes
- JSON/JSONB support for flexible fields
- Robust backup and recovery options

#### MongoDB (Order Service)
**Why MongoDB?**
- Flexible schema for varied order structures
- Excellent for document-based data (orders, certificates)
- High write throughput for order processing
- Natural fit for cart data (nested items)
- Easy horizontal scaling

### Caching Strategy (Redis)

**Project Service Caching:**
```
Cache Key Pattern: "project:list:*"
TTL: 5 minutes
Invalidation: On project create/update/delete
```

**Benefits:**
- Reduced database load
- Faster response times for frequently accessed data
- Lower latency for project listings

### Search Strategy (Elasticsearch)

**Project Search Indexing:**
```json
{
  "mappings": {
    "properties": {
      "name": { "type": "text" },
      "description": { "type": "text" },
      "category": { "type": "keyword" },
      "region": { "type": "keyword" },
      "country": { "type": "keyword" }
    }
  }
}
```

**Benefits:**
- Full-text search across project fields
- Fuzzy matching for typos
- Faceted search (filters)
- Fast search results

## Communication Patterns

### 1. Synchronous Communication (HTTP/REST)

**Gateway to Services:**
```
Client → Gateway → Service
  ↓        ↓         ↓
 HTTP    Proxy     HTTP
```

**Characteristics:**
- Request-response pattern
- Timeout: 30 seconds
- Connection pooling
- Retry logic (at client level)

### 2. Asynchronous Communication (RabbitMQ)

**Certificate Generation:**
```
Order Service → RabbitMQ → Certificate Consumer
     │             │              │
  Publish      Queue          Consume
  Message      Store          Process
                              Generate
```

**Message Format:**
```json
{
  "order_id": "ORD-20240101-ABC123",
  "user_id": "1",
  "items": [...],
  "total_tons": 10
}
```

**Benefits:**
- Decouples order creation from certificate generation
- Resilient to service failures
- Scalable message processing
- Guaranteed delivery

## Security Architecture

### Authentication & Authorization

#### JWT Token Structure
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "user_id": 1,
    "role": "user",
    "exp": 1704067200
  },
  "signature": "..."
}
```

#### Token Flow
```
1. User logs in → Service generates JWT
2. Client receives token
3. Client includes token in Authorization header
4. Gateway validates token
5. Gateway forwards request with token
6. Service processes request
```

### Security Layers

```
┌─────────────────────────────────┐
│  Layer 1: Rate Limiting         │  ← Prevent abuse
├─────────────────────────────────┤
│  Layer 2: CORS                  │  ← Control origins
├─────────────────────────────────┤
│  Layer 3: JWT Validation        │  ← Authenticate users
├─────────────────────────────────┤
│  Layer 4: Role-based Access     │  ← Authorize actions
├─────────────────────────────────┤
│  Layer 5: Input Validation      │  ← Prevent injection
└─────────────────────────────────┘
```

### Security Best Practices

1. **JWT Secrets**: Strong, random, environment-specific
2. **Password Hashing**: bcrypt with appropriate cost factor
3. **HTTPS**: TLS/SSL in production
4. **Database Credentials**: Separate per service, rotated regularly
5. **Service Communication**: Internal network isolation
6. **Secrets Management**: Environment variables, not hardcoded

## Scalability & Performance

### Horizontal Scaling

```
Load Balancer
     │
     ├─── API Gateway Instance 1
     ├─── API Gateway Instance 2
     └─── API Gateway Instance 3
          │
          ├─── User Service Pool
          ├─── Project Service Pool
          └─── Order Service Pool
```

### Performance Optimizations

#### 1. Caching (Redis)
- **Hit Ratio Target**: >80%
- **Cache Warming**: On service startup
- **Eviction Policy**: LRU

#### 2. Database Optimization
- **Indexes**: On frequently queried columns
- **Connection Pooling**: Managed by GORM
- **Query Optimization**: N+1 prevention

#### 3. Gateway Optimization
- **Request Pooling**: Reuse HTTP connections
- **Timeout Management**: Prevent hanging requests
- **Circuit Breaker**: (Future enhancement)

### Monitoring & Observability

**Metrics to Track:**
- Request rate (req/sec)
- Response time (p50, p95, p99)
- Error rate (%)
- Service health
- Database connections
- Cache hit/miss ratio
- Queue depth

**Logging Strategy:**
```
Gateway: All requests + routing decisions
Services: Business logic events + errors
Databases: Slow queries
Message Queue: Message processing
```

## Deployment Architecture

### Development Environment
```
Local Machine
├── Docker Compose
│   ├── All services
│   ├── All databases
│   └── All supporting services
```

### Production Environment (Recommended)
```
Cloud Infrastructure (AWS/GCP/Azure)
├── Load Balancer (ALB/ELB)
├── Container Orchestration (Kubernetes)
│   ├── API Gateway Pods (3+)
│   ├── Service Pods (2+ each)
│   └── Auto-scaling enabled
├── Managed Databases
│   ├── RDS PostgreSQL (Multi-AZ)
│   ├── DocumentDB/MongoDB Atlas
│   └── ElastiCache Redis
└── Managed Message Queue
    └── Amazon MQ / CloudAMQP
```

## Future Enhancements

### Short Term
- [ ] Circuit breaker pattern in Gateway
- [ ] Request tracing (distributed tracing)
- [ ] Metrics collection (Prometheus)
- [ ] Service mesh (Istio)

### Long Term
- [ ] GraphQL Gateway option
- [ ] Event sourcing for orders
- [ ] CQRS pattern implementation
- [ ] Multi-region deployment
- [ ] Advanced caching strategies

## Conclusion

The Carbon Clear architecture is designed for:
- **Scalability**: Horizontal scaling at each layer
- **Resilience**: Failure isolation between services
- **Performance**: Caching and async processing
- **Security**: Multiple security layers
- **Maintainability**: Clear separation of concerns

For deployment instructions, see [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md).
For testing procedures, see [TESTING_GUIDE.md](TESTING_GUIDE.md).

