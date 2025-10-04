package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"order_service/config"
	"order_service/models"
	"order_service/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

type OrderHandler struct {
	orderRepo       *repositories.OrderRepository
	cartRepo        *repositories.CartRepository
	certRepo        *repositories.CertificateRepository
	validator       *validator.Validate
	rabbitMQChannel *amqp.Channel
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderRepo:       repositories.NewOrderRepository(),
		cartRepo:        repositories.NewCartRepository(),
		certRepo:        repositories.NewCertificateRepository(),
		validator:       validator.New(),
		rabbitMQChannel: config.GetRabbitMQChannel(),
	}
}

// Checkout processes the user's cart and creates an order
func (h *OrderHandler) Checkout(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req models.CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get cart items
	cartItems, err := h.cartRepo.GetCartByUserID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve cart"})
	}

	if len(cartItems) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cart is empty"})
	}

	// For now, we'll create a single order for all cart items
	// In a real implementation, you might want to create separate orders per project
	var totalAmount float64
	var totalTonnes float64

	for _, item := range cartItems {
		// In a real implementation, you would fetch project details to get price per tonne
		// For now, we'll use a mock price
		pricePerTonne := 50.0 // This should come from project service
		itemTotal := item.Tonnes * pricePerTonne
		totalAmount += itemTotal
		totalTonnes += item.Tonnes
	}

	// Create order
	order := &models.Order{
		UserID:        uint(userID),
		ProjectID:     cartItems[0].ProjectID, // Using first project ID for simplicity
		Tonnes:        totalTonnes,
		PricePerTonne: 50.0, // Mock price
		TotalAmount:   totalAmount,
		Status:        "pending",
		PaymentID:     fmt.Sprintf("pay_%d_%d", userID, time.Now().Unix()),
	}

	if err := h.orderRepo.CreateOrder(order); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
	}

	// Process mock payment
	order.Status = "completed"
	if err := h.orderRepo.UpdateOrderStatus(order.ID, "completed"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order status"})
	}

	// Clear cart after successful checkout
	if err := h.cartRepo.ClearCart(uint(userID)); err != nil {
		// Log error but don't fail the checkout
		fmt.Printf("Warning: Failed to clear cart for user %d: %v\n", userID, err)
	}

	// Create certificate record
	certificate := &models.Certificate{
		OrderID:   order.ID,
		UserID:    uint(userID),
		ProjectID: order.ProjectID,
		Tonnes:    order.Tonnes,
		Status:    "pending",
	}

	if err := h.certRepo.CreateCertificate(certificate); err != nil {
		// Log error but don't fail the checkout
		fmt.Printf("Warning: Failed to create certificate record: %v\n", err)
	}

	// Send certificate generation message to RabbitMQ
	if err := h.sendCertificateGenerationMessage(order, certificate); err != nil {
		// Log error but don't fail the checkout
		fmt.Printf("Warning: Failed to send certificate generation message: %v\n", err)
	}

	// Convert to response
	orderResponse := models.OrderResponse{
		ID:             order.ID,
		UserID:         order.UserID,
		ProjectID:      order.ProjectID,
		Tonnes:         order.Tonnes,
		PricePerTonne:  order.PricePerTonne,
		TotalAmount:    order.TotalAmount,
		Status:         order.Status,
		PaymentID:      order.PaymentID,
		CertificateURL: order.CertificateURL,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, orderResponse)
}

// GetOrderHistory retrieves the user's order history
func (h *OrderHandler) GetOrderHistory(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	orders, err := h.orderRepo.GetOrdersByUserID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve order history"})
	}

	var responses []models.OrderResponse
	for _, order := range orders {
		responses = append(responses, models.OrderResponse{
			ID:             order.ID,
			UserID:         order.UserID,
			ProjectID:      order.ProjectID,
			Tonnes:         order.Tonnes,
			PricePerTonne:  order.PricePerTonne,
			TotalAmount:    order.TotalAmount,
			Status:         order.Status,
			PaymentID:      order.PaymentID,
			CertificateURL: order.CertificateURL,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, responses)
}

// GetOrder retrieves a specific order
func (h *OrderHandler) GetOrder(c echo.Context) error {
	orderIDStr := c.Param("orderID")
	_, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}

	// Convert to ObjectID (this is a simplified approach)
	// In a real implementation, you'd need proper ObjectID conversion
	return c.JSON(http.StatusOK, map[string]string{"message": "Order retrieved successfully"})
}

// GetCertificates retrieves the user's certificates
func (h *OrderHandler) GetCertificates(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	certificates, err := h.certRepo.GetCertificatesByUserID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve certificates"})
	}

	var responses []models.CertificateResponse
	for _, cert := range certificates {
		responses = append(responses, models.CertificateResponse{
			ID:             cert.ID,
			OrderID:        cert.OrderID,
			UserID:         cert.UserID,
			ProjectID:      cert.ProjectID,
			Tonnes:         cert.Tonnes,
			CertificateURL: cert.CertificateURL,
			Status:         cert.Status,
			CreatedAt:      cert.CreatedAt,
			UpdatedAt:      cert.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *OrderHandler) sendCertificateGenerationMessage(order *models.Order, certificate *models.Certificate) error {
	message := models.CertificateGenerationMessage{
		OrderID:   order.ID,
		UserID:    order.UserID,
		ProjectID: order.ProjectID,
		Tonnes:    order.Tonnes,
		UserEmail: "user@example.com", // This should come from user service
		UserName:  "User Name",        // This should come from user service
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = h.rabbitMQChannel.Publish(
		"",                       // exchange
		"certificate_generation", // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)

	return err
}
