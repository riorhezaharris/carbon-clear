package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

// CustomLogger provides detailed request logging
func CustomLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Log request details
			req := c.Request()
			res := c.Response()

			latency := time.Since(start)

			fmt.Printf("[%s] %s %s | Status: %d | Latency: %v | IP: %s | User-Agent: %s\n",
				start.Format("2006-01-02 15:04:05"),
				req.Method,
				req.URL.Path,
				res.Status,
				latency,
				c.RealIP(),
				req.UserAgent(),
			)

			return err
		}
	}
}
