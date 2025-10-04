package main

import (
	"log"
	"os"

	"order_service/config"
	"order_service/routes"
	"order_service/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to MongoDB
	config.ConnectMongoDB()
	defer config.CloseMongoDB()

	// Connect to RabbitMQ
	config.ConnectRabbitMQ()
	defer config.CloseRabbitMQ()

	// Initialize services
	certificateService := services.NewCertificateService()
	schedulerService := services.NewSchedulerService()

	// Start certificate consumer
	go certificateService.StartCertificateConsumer()

	// Start scheduler
	schedulerService.StartScheduler()
	defer schedulerService.StopScheduler()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	routes.SetupRoutes(e)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Order service starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
