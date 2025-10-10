package middleware

import (
	"api_gateway/config"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// UserJWTMiddleware validates user JWT tokens
func UserJWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.UserJWTSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "Unauthorized",
				"message": "Invalid or missing user authentication token",
			})
		},
	})
}

// AdminJWTMiddleware validates admin JWT tokens
func AdminJWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.AdminJWTSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "Unauthorized",
				"message": "Invalid or missing admin authentication token",
			})
		},
	})
}

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*JwtCustomClaims)

			if claims.Role != requiredRole {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error":   "Forbidden",
					"message": "You don't have permission to access this resource",
				})
			}

			return next(c)
		}
	}
}
