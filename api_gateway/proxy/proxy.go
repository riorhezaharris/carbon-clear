package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ServiceProxy struct {
	client *http.Client
}

func NewServiceProxy() *ServiceProxy {
	return &ServiceProxy{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyRequest forwards the request to the target service
func (sp *ServiceProxy) ProxyRequest(targetURL string) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		// Construct the full target URL
		fullURL := targetURL + req.URL.Path
		if req.URL.RawQuery != "" {
			fullURL += "?" + req.URL.RawQuery
		}

		// Read the request body
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create a new request
		proxyReq, err := http.NewRequest(req.Method, fullURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "Gateway Error",
				"message": "Failed to create proxy request",
			})
		}

		// Copy headers from original request
		for key, values := range req.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Execute the request
		resp, err := sp.client.Do(proxyReq)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{
				"error":   "Service Unavailable",
				"message": fmt.Sprintf("Failed to connect to backend service: %v", err),
			})
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "Gateway Error",
				"message": "Failed to read response from backend service",
			})
		}

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response().Header().Add(key, value)
			}
		}

		// Return the response
		return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
	}
}
