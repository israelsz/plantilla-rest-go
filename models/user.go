package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id"   bson:"_id,omitempty"`
	Email     string             `json:"email"  bson:"email,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Rol       string             `json:"rol"   bson:"rol"`
	Hash      string             `json:"_hash" bson:"_hash,omitempty"`
	Password  string             `json:"password,omitempty"  bson:"password,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}