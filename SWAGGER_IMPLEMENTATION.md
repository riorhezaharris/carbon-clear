# Swagger API Documentation Implementation

This document provides an overview of the Swagger/OpenAPI documentation implementation across all Carbon Clear microservices.

## Overview

Swagger documentation has been successfully implemented for all three microservices:
- **User Service** (Port 8080)
- **Project Service** (Port 8081)  
- **Order Service** (Port 8080)

## Accessing Swagger UI

### User Service
- **URL**: http://localhost:8080/swagger/index.html
- **Port**: 8080
- **Base Path**: `/`

### Project Service
- **URL**: http://localhost:8081/swagger/index.html
- **Port**: 8081
- **Base Path**: `/`

### Order Service
- **URL**: http://localhost:8080/swagger/index.html
- **Port**: 8080 (Note: Different from User Service, should be run separately or on different port)
- **Base Path**: `/`

## Features

All services include:

✅ **Interactive API Testing** - Test endpoints directly from the browser
✅ **Request/Response Examples** - View example requests and responses
✅ **Authentication Support** - JWT token authentication for protected endpoints
✅ **Model Schemas** - Detailed request and response models
✅ **Organized by Tags** - Endpoints grouped by functionality

## API Endpoints by Service

### User Service

#### Users Tag
- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - User login
- `GET /api/users/profile` - Get user profile (requires UserAuth)

#### Admin Tag
- `POST /admin/users/register` - Register a new admin
- `POST /admin/users/login` - Admin login
- `GET /admin/users` - Get all users (requires AdminAuth)
- `GET /admin/users/{id}` - Get user by ID (requires AdminAuth)
- `PUT /admin/users/{id}` - Update user (requires AdminAuth)
- `DELETE /admin/users/{id}` - Delete user (requires AdminAuth)

### Project Service

#### Projects Tag
- `GET /api/v1/projects` - Get all projects (with pagination)
- `GET /api/v1/projects/{id}` - Get project by ID
- `POST /api/v1/projects/search` - Search projects with filters
- `GET /api/v1/projects/categories` - Get available categories
- `GET /api/v1/projects/regions` - Get available regions
- `GET /api/v1/projects/countries` - Get available countries
- `POST /api/v1/projects/admin` - Create project (requires AdminAuth)
- `PUT /api/v1/projects/admin/{id}` - Update project (requires AdminAuth)
- `DELETE /api/v1/projects/admin/{id}` - Delete project (requires AdminAuth)

### Order Service

#### Cart Tag
- `POST /api/v1/cart/{userID}/items` - Add item to cart
- `GET /api/v1/cart/{userID}` - Get cart items
- `PUT /api/v1/cart/{userID}/items/{projectID}` - Update cart item
- `DELETE /api/v1/cart/{userID}/items/{projectID}` - Remove item from cart
- `DELETE /api/v1/cart/{userID}` - Clear cart

#### Orders Tag
- `POST /api/v1/orders/{userID}/checkout` - Checkout cart
- `GET /api/v1/orders/{userID}/history` - Get order history
- `GET /api/v1/orders/{orderID}` - Get order by ID
- `GET /api/v1/orders/{userID}/certificates` - Get user certificates

#### Admin Tag
- `GET /api/v1/admin/reports/monthly` - Get monthly report (requires AdminAuth)
- `GET /api/v1/admin/orders/date-range` - Get orders by date range (requires AdminAuth)
- `GET /api/v1/admin/statistics` - Get order statistics (requires AdminAuth)

## Authentication

### User Service
- **UserAuth**: JWT token for user endpoints (Header: `Authorization: Bearer <token>`)
- **AdminAuth**: JWT token for admin endpoints (Header: `Authorization: Bearer <token>`)

### Project Service
- **AdminAuth**: JWT token for admin endpoints (Header: `Authorization: Bearer <token>`)

### Order Service
- **UserAuth**: JWT token for user endpoints (Header: `Authorization: Bearer <token>`)
- **AdminAuth**: JWT token for admin endpoints (Header: `Authorization: Bearer <token>`)

## Implementation Details

### Dependencies Added
Each service now includes:
- `github.com/swaggo/swag` - Swagger documentation generator
- `github.com/swaggo/echo-swagger` - Echo framework Swagger middleware

### Files Generated
For each service, the following files are generated in the `/docs` folder:
- `docs.go` - Go documentation file
- `swagger.json` - JSON specification
- `swagger.yaml` - YAML specification

### Annotations Used
- `@title` - API title
- `@version` - API version
- `@description` - API description
- `@host` - API host
- `@BasePath` - API base path
- `@securityDefinitions.apikey` - JWT authentication definition
- `@Summary` - Endpoint summary
- `@Description` - Endpoint description
- `@Tags` - Endpoint grouping
- `@Accept` - Content type accepted
- `@Produce` - Content type produced
- `@Param` - Request parameters
- `@Success` - Success response
- `@Failure` - Error response
- `@Router` - Route definition
- `@Security` - Security requirement

## Regenerating Documentation

If you make changes to any API handler, regenerate the Swagger documentation:

```bash
# Navigate to the service directory
cd <service_directory>

# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g main.go
```

Or use the swag binary directly:

```bash
~/go/bin/swag init -g main.go
```

## Testing the Documentation

1. Start the service:
```bash
docker-compose up -d
# or
go run main.go
```

2. Open browser and navigate to the Swagger UI URL for the service

3. Try out the endpoints using the "Try it out" button

4. For protected endpoints:
   - First login via the login endpoint
   - Copy the token from the response
   - Click "Authorize" button in Swagger UI
   - Enter: `Bearer <your-token>`
   - Click "Authorize"
   - Now you can test protected endpoints

## Notes

- All services are fully documented with comprehensive Swagger annotations
- Each endpoint includes request/response examples
- Authentication is properly configured for protected routes
- Model definitions are automatically generated from Go structs
- Documentation is automatically updated when running `swag init`

## Troubleshooting

### Documentation not showing up
1. Ensure the service is running
2. Check that `/docs` folder exists with generated files
3. Verify the Swagger route is registered in main.go
4. Check that docs are imported: `_ "service_name/docs"`

### Authentication not working
1. Ensure you've obtained a valid JWT token via login
2. Format the token correctly: `Bearer <token>`
3. Use the "Authorize" button in Swagger UI
4. Verify the correct security definition is used for the endpoint

### Changes not reflected
1. Run `swag init -g main.go` after making changes
2. Restart the service
3. Hard refresh the browser (Ctrl+F5 / Cmd+Shift+R)

## Future Enhancements

Possible improvements for the Swagger documentation:
- Add more detailed examples for complex request bodies
- Include authentication flows in the documentation
- Add response headers documentation
- Include rate limiting information
- Add deprecation notices for older endpoints

