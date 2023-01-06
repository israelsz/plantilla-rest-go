package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cat struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Breed     string             `json:"breed" bson:"breed,omitempty"`
	Age       int                `json:"age" bson:"age,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}
