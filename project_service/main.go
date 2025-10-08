package main

import (
	"log"
	"net/http"
	"os"
	"project_service/config"
	_ "project_service/docs"
	"project_service/handlers"
	"project_service/repositories"
	"project_service/routes"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Project Service API
// @version 1.0
// @description This is the Project Marketplace Service API for Carbon Clear application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@carbonclear.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization
// @description Admin JWT token (format: Bearer <token>)

func main() {
	// Initialize Echo
	e := echo.New()

	// Health check endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Project Marketplace Service is running",
			"service": "project_service",
			"version": "1.0.0",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "project_service",
		})
	})

	// Initialize database
	_, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize Elasticsearch (optional)
	esClient, err := config.InitElasticsearch()
	if err != nil {
		log.Printf("Warning: Failed to connect to Elasticsearch: %v", err)
		log.Println("Continuing without Elasticsearch...")
	} else {
		log.Println("Successfully connected to Elasticsearch")
	}

	// Initialize Redis (optional)
	redisClient, err := config.InitRedis()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Continuing without Redis caching...")
	} else {
		log.Println("Successfully connected to Redis")
	}

	// Initialize repository and handler with external services
	projectRepo := repositories.NewProjectRepository()
	if esClient != nil {
		projectRepo.SetElasticsearchClient(esClient)
	}

	projectHandler := handlers.NewProjectHandler()
	if redisClient != nil {
		projectHandler.SetRedisClient(redisClient)
	}

	// Set up routes
	routes.ProjectRoute(e)

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Different port from user_service
	}

	log.Printf("Project Marketplace Service starting on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
