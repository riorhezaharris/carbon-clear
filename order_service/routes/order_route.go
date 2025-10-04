package routes

import (
	"order_service/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Initialize handlers
	cartHandler := handlers.NewCartHandler()
	orderHandler := handlers.NewOrderHandler()
	adminHandler := handlers.NewAdminHandler()

	// API version group
	api := e.Group("/api/v1")

	// Cart routes
	cart := api.Group("/cart")
	cart.POST("/:userID/items", cartHandler.AddToCart)
	cart.GET("/:userID", cartHandler.GetCart)
	cart.PUT("/:userID/items/:projectID", cartHandler.UpdateCartItem)
	cart.DELETE("/:userID/items/:projectID", cartHandler.RemoveFromCart)
	cart.DELETE("/:userID", cartHandler.ClearCart)

	// Order routes
	orders := api.Group("/orders")
	orders.POST("/:userID/checkout", orderHandler.Checkout)
	orders.GET("/:userID/history", orderHandler.GetOrderHistory)
	orders.GET("/:orderID", orderHandler.GetOrder)
	orders.GET("/:userID/certificates", orderHandler.GetCertificates)

	// Admin routes
	admin := api.Group("/admin")
	admin.GET("/reports/monthly", adminHandler.GetMonthlyReport)
	admin.GET("/orders/date-range", adminHandler.GetOrdersByDateRange)
	admin.GET("/statistics", adminHandler.GetOrderStatistics)
}
