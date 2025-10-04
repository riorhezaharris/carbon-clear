package repositories

import (
	"context"
	"time"

	"order_service/config"
	"order_service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository() *OrderRepository {
	db := config.GetMongoDB()
	return &OrderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *OrderRepository) GetOrderByID(id primitive.ObjectID) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order models.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(id primitive.ObjectID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func (r *OrderRepository) UpdateOrderCertificateURL(id primitive.ObjectID, certificateURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"certificate_url": certificateURL,
				"updated_at":      time.Now(),
			},
		},
	)
	return err
}

func (r *OrderRepository) GetOrdersByDateRange(startDate, endDate time.Time) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetMonthlyReport(year int, month int) (*models.MonthlyReport, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	orders, err := r.GetOrdersByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalOrders int
	var totalTonnes, totalRevenue float64

	for _, order := range orders {
		if order.Status == "completed" {
			totalOrders++
			totalTonnes += order.Tonnes
			totalRevenue += order.TotalAmount
		}
	}

	report := &models.MonthlyReport{
		Month:        startDate.Format("January"),
		Year:         year,
		TotalOrders:  totalOrders,
		TotalTonnes:  totalTonnes,
		TotalRevenue: totalRevenue,
		Orders:       convertToOrderResponses(orders),
	}

	return report, nil
}

func convertToOrderResponses(orders []models.Order) []models.OrderResponse {
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
	return responses
}
