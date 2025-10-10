#!/bin/bash

# Carbon Clear - Quick Start Script
# This script helps you quickly set up and start all services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  Carbon Clear - Quick Start${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}! $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    print_success "Docker is installed"
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    print_success "Docker Compose is installed"
    
    echo ""
}

# Generate environment file
generate_env() {
    if [ ! -f .env ]; then
        print_info "Creating .env file..."
        
        # Generate random secrets
        USER_SECRET=$(openssl rand -base64 32 2>/dev/null || echo "user-secret-$(date +%s)")
        ADMIN_SECRET=$(openssl rand -base64 32 2>/dev/null || echo "admin-secret-$(date +%s)")
        
        cat > .env << EOF
# Carbon Clear - Environment Configuration
# Generated on $(date)

# JWT Secrets
USER_JWT_SECRET=${USER_SECRET}
ADMIN_JWT_SECRET=${ADMIN_SECRET}

# Database Passwords
POSTGRES_PASSWORD=postgres
MONGODB_PASSWORD=admin123
REDIS_PASSWORD=redis123

# RabbitMQ Configuration
RABBITMQ_USER=admin
RABBITMQ_PASSWORD=admin123
EOF
        
        print_success ".env file created"
        print_warning "Remember to change passwords in production!"
    else
        print_info ".env file already exists, skipping..."
    fi
    echo ""
}

# Build services
build_services() {
    print_info "Building Docker images..."
    docker-compose build
    print_success "Docker images built successfully"
    echo ""
}

# Start services
start_services() {
    print_info "Starting all services..."
    docker-compose up -d
    print_success "Services started"
    echo ""
}

# Wait for services to be healthy
wait_for_services() {
    print_info "Waiting for services to be healthy..."
    echo ""
    
    # Wait for databases
    print_info "Waiting for databases..."
    sleep 10
    
    # Check API Gateway
    print_info "Checking API Gateway..."
    for i in {1..30}; do
        if curl -f -s http://localhost:8000/health > /dev/null 2>&1; then
            print_success "API Gateway is healthy"
            break
        fi
        if [ $i -eq 30 ]; then
            print_error "API Gateway failed to start"
            exit 1
        fi
        sleep 2
    done
    
    # Check User Service
    print_info "Checking User Service..."
    for i in {1..30}; do
        if curl -f -s http://localhost:8082/ > /dev/null 2>&1; then
            print_success "User Service is healthy"
            break
        fi
        if [ $i -eq 30 ]; then
            print_error "User Service failed to start"
            exit 1
        fi
        sleep 2
    done
    
    # Check Project Service
    print_info "Checking Project Service..."
    for i in {1..30}; do
        if curl -f -s http://localhost:8081/health > /dev/null 2>&1; then
            print_success "Project Service is healthy"
            break
        fi
        if [ $i -eq 30 ]; then
            print_error "Project Service failed to start"
            exit 1
        fi
        sleep 2
    done
    
    # Check Order Service
    print_info "Checking Order Service..."
    for i in {1..30}; do
        if curl -f -s http://localhost:8080/swagger/ > /dev/null 2>&1; then
            print_success "Order Service is healthy"
            break
        fi
        if [ $i -eq 30 ]; then
            print_warning "Order Service may still be starting..."
        fi
        sleep 2
    done
    
    echo ""
}

# Print access information
print_access_info() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  All Services are Running!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "Access the services at:"
    echo ""
    echo -e "  ${BLUE}API Gateway:${NC}          http://localhost:8000"
    echo -e "  ${BLUE}API Documentation:${NC}    http://localhost:8000/swagger/index.html"
    echo ""
    echo -e "  ${BLUE}User Service:${NC}         http://localhost:8082"
    echo -e "  ${BLUE}User Service Docs:${NC}    http://localhost:8082/swagger/index.html"
    echo ""
    echo -e "  ${BLUE}Project Service:${NC}      http://localhost:8081"
    echo -e "  ${BLUE}Project Service Docs:${NC} http://localhost:8081/swagger/index.html"
    echo ""
    echo -e "  ${BLUE}Order Service:${NC}        http://localhost:8080"
    echo -e "  ${BLUE}Order Service Docs:${NC}   http://localhost:8080/swagger/index.html"
    echo ""
    echo -e "  ${BLUE}RabbitMQ Console:${NC}     http://localhost:15672"
    echo -e "    Username: ${YELLOW}admin${NC}"
    echo -e "    Password: ${YELLOW}admin123${NC}"
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "Useful commands:"
    echo ""
    echo "  View logs:           docker-compose logs -f"
    echo "  Stop services:       docker-compose down"
    echo "  Restart services:    docker-compose restart"
    echo "  Check status:        docker-compose ps"
    echo ""
    echo "  Or use Make commands:"
    echo "    make logs          View all logs"
    echo "    make down          Stop all services"
    echo "    make restart       Restart all services"
    echo "    make status        Check service status"
    echo "    make test          Run health checks"
    echo ""
}

# Main execution
main() {
    print_header
    check_prerequisites
    generate_env
    build_services
    start_services
    wait_for_services
    print_access_info
}

# Run main function
main

