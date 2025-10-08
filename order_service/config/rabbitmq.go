package config

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var RabbitMQConn *amqp.Connection
var RabbitMQChannel *amqp.Channel

func ConnectRabbitMQ() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}

	// Get queue name from environment variable
	queueName := os.Getenv("CERTIFICATE_QUEUE_NAME")
	if queueName == "" {
		queueName = "certificate_generation"
	}

	// Declare the certificate generation queue
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare certificate generation queue:", err)
	}

	RabbitMQConn = conn
	RabbitMQChannel = ch

	fmt.Println("Connected to RabbitMQ successfully")
}

func GetRabbitMQChannel() *amqp.Channel {
	return RabbitMQChannel
}

func CloseRabbitMQ() {
	if RabbitMQChannel != nil {
		RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}
