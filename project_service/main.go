package main

import (
	"log"
	"net/http"
	"os"
	"project_service/config"
	"project_service/handlers"
	"project_service/repositories"
	"project_service/routes"

	"github.com/labstack/echo/v4"
)

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

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Different port from user_service
	}

	log.Printf("Project Marketplace Service starting on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
