package main

import (
	"api_gateway/config"
	_ "api_gateway/docs"
	"api_gateway/middleware"
	"api_gateway/proxy"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Carbon Clear API Gateway
// @version 1.0
// @description API Gateway for Carbon Clear microservices architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@carbonclear.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey UserAuth
// @in header
// @name Authorization
// @description User JWT token (format: Bearer <token>)

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization
// @description Admin JWT token (format: Bearer <token>)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Create Echo instance
	e := echo.New()

	// Basic middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.CustomLogger())

	// Rate limiting middleware
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitPerMin)
	e.Use(rateLimiter.Middleware())

	// Setup routes (proxy to backend services)
	proxy.SetupRoutes(e, cfg)

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	log.Printf("API Gateway starting on port %s", cfg.Port)
	log.Printf("User Service URL: %s", cfg.UserServiceURL)
	log.Printf("Project Service URL: %s", cfg.ProjectServiceURL)
	log.Printf("Order Service URL: %s", cfg.OrderServiceURL)
	log.Printf("Rate Limit: %d requests/minute", cfg.RateLimitPerMin)

	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}
