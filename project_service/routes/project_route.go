package routes

import (
	"os"
	"project_service/handlers"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type JwtCustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func ProjectRoute(e *echo.Echo) {
	// Create handler instance
	projectHandler := handlers.NewProjectHandler()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API group
	api := e.Group("/api/v1")

	// Project routes
	projects := api.Group("/projects")

	// Public routes (no authentication required)
	projects.GET("", projectHandler.GetAllProjects)                  // Browse Projects
	projects.GET("/:id", projectHandler.GetProject)                  // View Project Details
	projects.POST("/search", projectHandler.SearchProjects)          // Filter & Search Projects
	projects.GET("/categories", projectHandler.GetProjectCategories) // Get available categories
	projects.GET("/regions", projectHandler.GetProjectRegions)       // Get available regions
	projects.GET("/countries", projectHandler.GetProjectCountries)   // Get available countries

	// Admin routes (authentication required)
	admin := projects.Group("/admin")
	admin.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("ADMIN_JWT_SECRET")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
	}))

	admin.POST("", projectHandler.CreateProject)       // Create new project
	admin.PUT("/:id", projectHandler.UpdateProject)    // Update project
	admin.DELETE("/:id", projectHandler.DeleteProject) // Delete project
}
