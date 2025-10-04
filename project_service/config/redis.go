package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RedisClient *redis.Client

func InitRedis() (*redis.Client, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Redis configuration
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // default Redis address
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // default database

	// Parse Redis DB if provided
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			redisDB = db
		}
	}

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %v", err)
	}

	RedisClient = client
	log.Println("Successfully connected to Redis")
	return client, nil
}

// Cache operations
func SetCache(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil // Skip if Redis is not configured
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling cache value: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = RedisClient.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("error setting cache: %v", err)
	}

	return nil
}

func GetCache(key string, dest interface{}) error {
	if RedisClient == nil {
		return redis.Nil // Skip if Redis is not configured
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("error unmarshaling cache value: %v", err)
	}

	return nil
}

func DeleteCache(key string) error {
	if RedisClient == nil {
		return nil // Skip if Redis is not configured
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting cache: %v", err)
	}

	return nil
}

func DeleteCachePattern(pattern string) error {
	if RedisClient == nil {
		return nil // Skip if Redis is not configured
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all keys matching the pattern
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("error getting keys: %v", err)
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all matching keys
	err = RedisClient.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("error deleting keys: %v", err)
	}

	return nil
}

// Cache key generators
func GetProjectCacheKey(id string) string {
	return fmt.Sprintf("project:%s", id)
}

func GetProjectsListCacheKey(limit, offset int) string {
	return fmt.Sprintf("projects:all:%d:%d", limit, offset)
}

func GetSearchCacheKey(query string, category, region, country []string, minPrice, maxPrice float64, limit, offset int) string {
	return fmt.Sprintf("search:%s:%v:%v:%v:%.2f:%.2f:%d:%d",
		query, category, region, country, minPrice, maxPrice, limit, offset)
}

func GetCategoriesCacheKey() string {
	return "project:categories"
}

func GetRegionsCacheKey() string {
	return "project:regions"
}

func GetCountriesCacheKey() string {
	return "project:countries"
}
