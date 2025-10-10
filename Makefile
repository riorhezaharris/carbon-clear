.PHONY: help build up down logs clean restart status test

# Default target
help:
	@echo "Carbon Clear - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  help        - Show this help message"
	@echo "  build       - Build all Docker images"
	@echo "  up          - Start all services"
	@echo "  down        - Stop all services"
	@echo "  logs        - View logs from all services"
	@echo "  clean       - Remove all containers, volumes, and images"
	@echo "  restart     - Restart all services"
	@echo "  status      - Show status of all services"
	@echo "  test        - Run health checks on all services"
	@echo ""
	@echo "Individual Service Commands:"
	@echo "  gateway-logs     - View API Gateway logs"
	@echo "  user-logs        - View User Service logs"
	@echo "  project-logs     - View Project Service logs"
	@echo "  order-logs       - View Order Service logs"
	@echo "  gateway-restart  - Restart API Gateway"
	@echo "  user-restart     - Restart User Service"
	@echo "  project-restart  - Restart Project Service"
	@echo "  order-restart    - Restart Order Service"

# Build all Docker images
build:
	@echo "Building all Docker images..."
	docker-compose build

# Start all services
up:
	@echo "Starting all services..."
	docker-compose up -d
	@echo ""
	@echo "Services started! Access points:"
	@echo "  API Gateway:      http://localhost:8000"
	@echo "  API Gateway Docs: http://localhost:8000/swagger/index.html"
	@echo "  User Service:     http://localhost:8082"
	@echo "  Project Service:  http://localhost:8081"
	@echo "  Order Service:    http://localhost:8080"
	@echo "  RabbitMQ Admin:   http://localhost:15672 (admin/admin123)"

# Stop all services
down:
	@echo "Stopping all services..."
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Clean everything
clean:
	@echo "Cleaning up..."
	docker-compose down -v --rmi all
	@echo "Cleanup complete!"

# Restart all services
restart:
	@echo "Restarting all services..."
	docker-compose restart

# Show status
status:
	@echo "Service Status:"
	docker-compose ps

# Run health checks
test:
	@echo "Running health checks..."
	@echo ""
	@echo "API Gateway:"
	@curl -f http://localhost:8000/health || echo "  ❌ Failed"
	@echo ""
	@echo "User Service:"
	@curl -f http://localhost:8082/ || echo "  ❌ Failed"
	@echo ""
	@echo "Project Service:"
	@curl -f http://localhost:8081/health || echo "  ❌ Failed"
	@echo ""
	@echo "Order Service:"
	@curl -f http://localhost:8080/swagger/ || echo "  ❌ Failed"

# Individual service commands
gateway-logs:
	docker-compose logs -f api_gateway

user-logs:
	docker-compose logs -f user_service

project-logs:
	docker-compose logs -f project_service

order-logs:
	docker-compose logs -f order_service

gateway-restart:
	docker-compose restart api_gateway

user-restart:
	docker-compose restart user_service

project-restart:
	docker-compose restart project_service

order-restart:
	docker-compose restart order_service

# Development commands
dev-gateway:
	cd api_gateway && go run main.go

dev-user:
	cd user_service && go run main.go

dev-project:
	cd project_service && go run main.go

dev-order:
	cd order_service && go run main.go

