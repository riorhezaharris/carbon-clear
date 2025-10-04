package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         uint               `json:"user_id" bson:"user_id"`
	ProjectID      uint               `json:"project_id" bson:"project_id"`
	Tonnes         float64            `json:"tonnes" bson:"tonnes"`
	PricePerTonne  float64            `json:"price_per_tonne" bson:"price_per_tonne"`
	TotalAmount    float64            `json:"total_amount" bson:"total_amount"`
	Status         string             `json:"status" bson:"status"` // pending, completed, cancelled
	PaymentID      string             `json:"payment_id" bson:"payment_id"`
	CertificateURL string             `json:"certificate_url" bson:"certificate_url"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type CartItem struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    uint               `json:"user_id" bson:"user_id"`
	ProjectID uint               `json:"project_id" bson:"project_id"`
	Tonnes    float64            `json:"tonnes" bson:"tonnes"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Certificate struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderID        primitive.ObjectID `json:"order_id" bson:"order_id"`
	UserID         uint               `json:"user_id" bson:"user_id"`
	ProjectID      uint               `json:"project_id" bson:"project_id"`
	Tonnes         float64            `json:"tonnes" bson:"tonnes"`
	CertificateURL string             `json:"certificate_url" bson:"certificate_url"`
	Status         string             `json:"status" bson:"status"` // pending, generated, failed
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

// Request/Response DTOs
type AddToCartRequest struct {
	ProjectID uint    `json:"project_id" validate:"required"`
	Tonnes    float64 `json:"tonnes" validate:"required,gt=0"`
}

type CheckoutRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
}

type OrderResponse struct {
	ID             primitive.ObjectID `json:"id"`
	UserID         uint               `json:"user_id"`
	ProjectID      uint               `json:"project_id"`
	Tonnes         float64            `json:"tonnes"`
	PricePerTonne  float64            `json:"price_per_tonne"`
	TotalAmount    float64            `json:"total_amount"`
	Status         string             `json:"status"`
	PaymentID      string             `json:"payment_id"`
	CertificateURL string             `json:"certificate_url"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type CartItemResponse struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    uint               `json:"user_id"`
	ProjectID uint               `json:"project_id"`
	Tonnes    float64            `json:"tonnes"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type CertificateResponse struct {
	ID             primitive.ObjectID `json:"id"`
	OrderID        primitive.ObjectID `json:"order_id"`
	UserID         uint               `json:"user_id"`
	ProjectID      uint               `json:"project_id"`
	Tonnes         float64            `json:"tonnes"`
	CertificateURL string             `json:"certificate_url"`
	Status         string             `json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type MonthlyReport struct {
	Month        string          `json:"month"`
	Year         int             `json:"year"`
	TotalOrders  int             `json:"total_orders"`
	TotalTonnes  float64         `json:"total_tonnes"`
	TotalRevenue float64         `json:"total_revenue"`
	Orders       []OrderResponse `json:"orders"`
}

type CertificateGenerationMessage struct {
	OrderID   primitive.ObjectID `json:"order_id"`
	UserID    uint               `json:"user_id"`
	ProjectID uint               `json:"project_id"`
	Tonnes    float64            `json:"tonnes"`
	UserEmail string             `json:"user_email"`
	UserName  string             `json:"user_name"`
}
