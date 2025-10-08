package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"order_service/config"
	"order_service/models"
	"order_service/repositories"

	"github.com/streadway/amqp"
)

type CertificateService struct {
	certRepo        *repositories.CertificateRepository
	orderRepo       *repositories.OrderRepository
	rabbitMQChannel *amqp.Channel
}

func NewCertificateService() *CertificateService {
	return &CertificateService{
		certRepo:        repositories.NewCertificateRepository(),
		orderRepo:       repositories.NewOrderRepository(),
		rabbitMQChannel: config.GetRabbitMQChannel(),
	}
}

// StartCertificateConsumer starts listening for certificate generation messages
func (s *CertificateService) StartCertificateConsumer() {
	// Get queue name from environment variable
	queueName := os.Getenv("CERTIFICATE_QUEUE_NAME")
	if queueName == "" {
		queueName = "certificate_generation"
	}

	msgs, err := s.rabbitMQChannel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	go func() {
		for msg := range msgs {
			var certMessage models.CertificateGenerationMessage
			if err := json.Unmarshal(msg.Body, &certMessage); err != nil {
				log.Printf("Failed to unmarshal certificate message: %v", err)
				continue
			}

			if err := s.generateCertificate(&certMessage); err != nil {
				log.Printf("Failed to generate certificate: %v", err)
			}
		}
	}()

	log.Println("Certificate consumer started")
}

func (s *CertificateService) generateCertificate(message *models.CertificateGenerationMessage) error {
	log.Printf("Generating certificate for order %s", message.OrderID.Hex())

	// Simulate certificate generation process
	time.Sleep(2 * time.Second)

	// Generate a mock certificate URL
	certificateURL := fmt.Sprintf("https://certificates.carbonclear.com/cert_%s.pdf", message.OrderID.Hex())

	// Update certificate record
	if err := s.certRepo.UpdateCertificateURL(message.OrderID, certificateURL); err != nil {
		return fmt.Errorf("failed to update certificate URL: %v", err)
	}

	// Update order with certificate URL
	if err := s.orderRepo.UpdateOrderCertificateURL(message.OrderID, certificateURL); err != nil {
		return fmt.Errorf("failed to update order certificate URL: %v", err)
	}

	// In a real implementation, you would:
	// 1. Generate an actual PDF certificate
	// 2. Upload it to a file storage service (S3, etc.)
	// 3. Send an email notification to the user
	// 4. Update the user's profile with the certificate

	log.Printf("Certificate generated successfully for order %s", message.OrderID.Hex())
	return nil
}

// GenerateMockPDF creates a simple mock PDF certificate
func (s *CertificateService) GenerateMockPDF(orderID string, userEmail, userName string, tonnes float64) ([]byte, error) {
	// This is a very basic mock implementation
	// In a real application, you would use a proper PDF library like gofpdf or unidoc

	content := fmt.Sprintf(`
CARBON OFFSET CERTIFICATE

Certificate ID: %s
Issued To: %s (%s)
Date: %s
Carbon Offset: %.2f tonnes CO2

This certificate confirms that the above individual has offset
%.2f tonnes of CO2 through our carbon offset program.

Thank you for your contribution to fighting climate change!

Carbon Clear Team
	`, orderID, userName, userEmail, time.Now().Format("January 2, 2006"), tonnes, tonnes)

	// In a real implementation, you would generate an actual PDF
	// For now, we'll return the content as bytes
	return []byte(content), nil
}

// SendCertificateEmail sends an email notification with the certificate
func (s *CertificateService) SendCertificateEmail(userEmail, userName, certificateURL string) error {
	// In a real implementation, you would:
	// 1. Use an email service like SendGrid, AWS SES, or similar
	// 2. Create an HTML email template
	// 3. Send the email with the certificate attachment or link

	log.Printf("Mock email sent to %s with certificate URL: %s", userEmail, certificateURL)
	return nil
}
