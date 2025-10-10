package config

import (
	"os"
)

type Config struct {
	Port              string
	UserServiceURL    string
	ProjectServiceURL string
	OrderServiceURL   string
	UserJWTSecret     string
	AdminJWTSecret    string
	RateLimitPerMin   int
}

func LoadConfig() *Config {
	return &Config{
		Port:              getEnv("PORT", "8000"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		ProjectServiceURL: getEnv("PROJECT_SERVICE_URL", "http://localhost:8081"),
		OrderServiceURL:   getEnv("ORDER_SERVICE_URL", "http://localhost:8080"),
		UserJWTSecret:     getEnv("USER_JWT_SECRET", "user-secret-key"),
		AdminJWTSecret:    getEnv("ADMIN_JWT_SECRET", "admin-secret-key"),
		RateLimitPerMin:   100,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
