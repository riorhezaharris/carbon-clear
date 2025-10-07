package config

import (
	"fmt"
	"log"
	"os"
	"project_service/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	var err error

	// Initiate secrets and credentials (optional - only for local development)
	err = godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables from system")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ping the database and error handling
	psql, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = psql.Ping()
	if err != nil {
		fmt.Println("Ping to database is failed")
		return nil, err
	}

	DB = db

	err = db.AutoMigrate(&models.Project{})
	if err != nil {
		return nil, err
	}

	return DB, nil
}
