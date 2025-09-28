package main

import (
	"log"
	"net/http"
	"user_service/configs"

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

	e.Logger.Fatal(e.Start(":8080"))
}
