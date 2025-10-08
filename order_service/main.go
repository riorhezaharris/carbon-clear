package main

import (
	"log"
	"os"

	"order_service/config"
	_ "order_service/docs"
	"order_service/routes"
	"order_service/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Order Service API
// @version 1.0
// @description This is the Order and Cart Service API for Carbon Clear application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@carbonclear.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
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

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
