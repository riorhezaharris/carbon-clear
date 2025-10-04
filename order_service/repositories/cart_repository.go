package repositories

import (
	"context"
	"time"

	"order_service/config"
	"order_service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository struct {
	collection *mongo.Collection
}

func NewCartRepository() *CartRepository {
	db := config.GetMongoDB()
	return &CartRepository{
		collection: db.Collection("cart_items"),
	}
}

func (r *CartRepository) AddToCart(cartItem *models.CartItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if item already exists for this user and project
	filter := bson.M{
		"user_id":    cartItem.UserID,
		"project_id": cartItem.ProjectID,
	}

	var existingItem models.CartItem
	err := r.collection.FindOne(ctx, filter).Decode(&existingItem)
	if err == nil {
		// Update existing item
		_, err = r.collection.UpdateOne(
			ctx,
			filter,
			bson.M{
				"$set": bson.M{
					"tonnes":     existingItem.Tonnes + cartItem.Tonnes,
					"updated_at": time.Now(),
				},
			},
		)
		return err
	} else if err == mongo.ErrNoDocuments {
		// Create new item
		cartItem.CreatedAt = time.Now()
		cartItem.UpdatedAt = time.Now()

		_, err = r.collection.InsertOne(ctx, cartItem)
		return err
	}

	return err
}

func (r *CartRepository) GetCartByUserID(userID uint) ([]models.CartItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cartItems []models.CartItem
	if err = cursor.All(ctx, &cartItems); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (r *CartRepository) RemoveFromCart(userID uint, projectID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(
		ctx,
		bson.M{
			"user_id":    userID,
			"project_id": projectID,
		},
	)
	return err
}

func (r *CartRepository) ClearCart(userID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}

func (r *CartRepository) UpdateCartItem(userID uint, projectID uint, tonnes float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tonnes <= 0 {
		return r.RemoveFromCart(userID, projectID)
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"user_id":    userID,
			"project_id": projectID,
		},
		bson.M{
			"$set": bson.M{
				"tonnes":     tonnes,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}
