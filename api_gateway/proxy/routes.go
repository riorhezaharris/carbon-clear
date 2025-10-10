package proxy

import (
	"api_gateway/config"
	"api_gateway/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, cfg *config.Config) {
	proxy := NewServiceProxy()

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "api_gateway",
			"version": "1.0.0",
		})
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Carbon Clear API Gateway",
			"version": "1.0.0",
			"services": map[string]string{
				"user_service":    cfg.UserServiceURL,
				"project_service": cfg.ProjectServiceURL,
				"order_service":   cfg.OrderServiceURL,
			},
		})
	})

	// API version group
	api := e.Group("/api")

	// ===== USER SERVICE ROUTES =====
	// Public user routes (no auth required)
	api.POST("/users/register", proxy.ProxyRequest(cfg.UserServiceURL))
	api.POST("/users/login", proxy.ProxyRequest(cfg.UserServiceURL))

	// Protected user routes (require user JWT)
	userProtected := api.Group("/users")
	userProtected.Use(middleware.UserJWTMiddleware(cfg))
	userProtected.GET("/profile", proxy.ProxyRequest(cfg.UserServiceURL))

	// Admin user management routes
	api.POST("/admin/users/register", proxy.ProxyRequest(cfg.UserServiceURL))
	api.POST("/admin/users/login", proxy.ProxyRequest(cfg.UserServiceURL))

	adminUsers := api.Group("/admin")
	adminUsers.Use(middleware.AdminJWTMiddleware(cfg))
	adminUsers.GET("/users", proxy.ProxyRequest(cfg.UserServiceURL))
	adminUsers.POST("/users", proxy.ProxyRequest(cfg.UserServiceURL))
	adminUsers.GET("/users/:id", proxy.ProxyRequest(cfg.UserServiceURL))
	adminUsers.PUT("/users/:id", proxy.ProxyRequest(cfg.UserServiceURL))
	adminUsers.DELETE("/users/:id", proxy.ProxyRequest(cfg.UserServiceURL))

	// ===== PROJECT SERVICE ROUTES =====
	// Public project routes (no auth required)
	projects := api.Group("/v1/projects")
	projects.GET("", proxy.ProxyRequest(cfg.ProjectServiceURL))
	projects.GET("/:id", proxy.ProxyRequest(cfg.ProjectServiceURL))
	projects.POST("/search", proxy.ProxyRequest(cfg.ProjectServiceURL))
	projects.GET("/categories", proxy.ProxyRequest(cfg.ProjectServiceURL))
	projects.GET("/regions", proxy.ProxyRequest(cfg.ProjectServiceURL))
	projects.GET("/countries", proxy.ProxyRequest(cfg.ProjectServiceURL))

	// Admin project routes (require admin JWT)
	adminProjects := projects.Group("/admin")
	adminProjects.Use(middleware.AdminJWTMiddleware(cfg))
	adminProjects.POST("", proxy.ProxyRequest(cfg.ProjectServiceURL))
	adminProjects.PUT("/:id", proxy.ProxyRequest(cfg.ProjectServiceURL))
	adminProjects.DELETE("/:id", proxy.ProxyRequest(cfg.ProjectServiceURL))

	// ===== ORDER SERVICE ROUTES =====
	// Cart routes (protected - require user JWT)
	cart := api.Group("/v1/cart")
	cart.Use(middleware.UserJWTMiddleware(cfg))
	cart.POST("/:userID/items", proxy.ProxyRequest(cfg.OrderServiceURL))
	cart.GET("/:userID", proxy.ProxyRequest(cfg.OrderServiceURL))
	cart.PUT("/:userID/items/:projectID", proxy.ProxyRequest(cfg.OrderServiceURL))
	cart.DELETE("/:userID/items/:projectID", proxy.ProxyRequest(cfg.OrderServiceURL))
	cart.DELETE("/:userID", proxy.ProxyRequest(cfg.OrderServiceURL))

	// Order routes (protected - require user JWT)
	orders := api.Group("/v1/orders")
	orders.Use(middleware.UserJWTMiddleware(cfg))
	orders.POST("/:userID/checkout", proxy.ProxyRequest(cfg.OrderServiceURL))
	orders.GET("/:userID/history", proxy.ProxyRequest(cfg.OrderServiceURL))
	orders.GET("/:orderID", proxy.ProxyRequest(cfg.OrderServiceURL))
	orders.GET("/:userID/certificates", proxy.ProxyRequest(cfg.OrderServiceURL))

	// Admin order routes (require admin JWT)
	adminOrders := api.Group("/v1/admin")
	adminOrders.Use(middleware.AdminJWTMiddleware(cfg))
	adminOrders.GET("/reports/monthly", proxy.ProxyRequest(cfg.OrderServiceURL))
	adminOrders.GET("/orders/date-range", proxy.ProxyRequest(cfg.OrderServiceURL))
	adminOrders.GET("/statistics", proxy.ProxyRequest(cfg.OrderServiceURL))
}
