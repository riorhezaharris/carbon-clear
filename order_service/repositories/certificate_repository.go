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

type CertificateRepository struct {
	collection *mongo.Collection
}

func NewCertificateRepository() *CertificateRepository {
	db := config.GetMongoDB()
	return &CertificateRepository{
		collection: db.Collection("certificates"),
	}
}

func (r *CertificateRepository) CreateCertificate(certificate *models.Certificate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	certificate.CreatedAt = time.Now()
	certificate.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, certificate)
	if err != nil {
		return err
	}

	certificate.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *CertificateRepository) GetCertificateByOrderID(orderID primitive.ObjectID) (*models.Certificate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var certificate models.Certificate
	err := r.collection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&certificate)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}

func (r *CertificateRepository) GetCertificatesByUserID(userID uint) ([]models.Certificate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certificates []models.Certificate
	if err = cursor.All(ctx, &certificates); err != nil {
		return nil, err
	}

	return certificates, nil
}

func (r *CertificateRepository) UpdateCertificateStatus(id primitive.ObjectID, status string) error {
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

func (r *CertificateRepository) UpdateCertificateURL(id primitive.ObjectID, certificateURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"certificate_url": certificateURL,
				"status":          "generated",
				"updated_at":      time.Now(),
			},
		},
	)
	return err
}
