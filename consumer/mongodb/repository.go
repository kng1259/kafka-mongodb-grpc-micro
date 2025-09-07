package mongodb

import (
	"context"
	"log"

	"kmgm-consumer/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}

func (r *Repository) StoreProduct(ctx context.Context, product models.Product) error {
	// Generate ID if not present
	if product.ID == primitive.NilObjectID {
		product.ID = primitive.NewObjectID()
	}

	// Store product in MongoDB
	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	log.Printf("Stored product in MongoDB: %s (ID: %s)", product.Name, product.ID.Hex())
	return nil
}

func (r *Repository) GetCollection() *mongo.Collection {
	return r.collection
}
