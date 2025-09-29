package routes

import (
	"os"
	"user_service/handlers"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func UserRoute(e *echo.Echo) {
	e.POST("/api/users/register", handlers.RegisterUser)
	e.POST("/api/users/login", handlers.LoginUser)

	e.POST("/admin/users/register", handlers.RegisterAdmin)
	e.POST("/admin/users/login", handlers.LoginAdmin)
	admin := e.Group("/admin")
	admin.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("ADMIN_JWT_SECRET")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
	}))
	admin.GET("/users", handlers.GetAllUsers)
	admin.POST("/users", handlers.RegisterUser)
	admin.GET("/users/:id", handlers.GetUserByID)
	admin.PUT("/users/:id", handlers.UpdateUser)
	admin.DELETE("/users/:id", handlers.DeleteUser)
}
