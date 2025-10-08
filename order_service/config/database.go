package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client
var MongoDBDatabase *mongo.Database

func ConnectMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	MongoDBClient = client

	// Get database name from environment variable
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "carbon_clear_orders"
	}
	MongoDBDatabase = client.Database(dbName)

	fmt.Println("Connected to MongoDB successfully")
}

func GetMongoDB() *mongo.Database {
	return MongoDBDatabase
}

func CloseMongoDB() {
	if MongoDBClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		MongoDBClient.Disconnect(ctx)
	}
}
