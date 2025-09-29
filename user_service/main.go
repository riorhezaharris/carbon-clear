package main

import (
	"log"
	"net/http"
	"os"
	"user_service/configs"
	"user_service/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	_, err := configs.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	routes.UserRoute(e)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
