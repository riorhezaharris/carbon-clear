package routes

import (
	"os"
	"user_service/handlers"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/api/users/register", handlers.RegisterUser)
	e.POST("/api/users/login", handlers.LoginUser)

	// User routes with JWT authentication
	user := e.Group("/api/users")
	user.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("USER_JWT_SECRET")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JwtCustomClaims)
		},
	}))
	user.GET("/profile", handlers.GetProfile)

	e.POST("/api/admin/users/register", handlers.RegisterAdmin)
	e.POST("/api/admin/users/login", handlers.LoginAdmin)
	admin := e.Group("/api/admin")
	admin.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("ADMIN_JWT_SECRET")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JwtCustomClaims)
		},
	}))
	admin.GET("/users", handlers.GetAllUsers)
	admin.POST("/users", handlers.RegisterUser)
	admin.GET("/users/:id", handlers.GetUserByID)
	admin.PUT("/users/:id", handlers.UpdateUser)
	admin.DELETE("/users/:id", handlers.DeleteUser)
}
