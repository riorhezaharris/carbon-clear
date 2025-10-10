package main

import (
	"log"
	"net/http"
	"os"
	"user_service/configs"
	_ "user_service/docs"
	"user_service/routes"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title User Service API
// @version 1.0
// @description This is the User Service API for Carbon Clear application
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
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Health Check!")
	})

	_, err := configs.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	routes.UserRoute(e)

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
