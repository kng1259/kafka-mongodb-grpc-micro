package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Price       float64            `json:"price" bson:"price" binding:"required"`
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
