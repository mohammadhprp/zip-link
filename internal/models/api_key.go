package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APIKey struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Key          string             `bson:"key" json:"key"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"`
	RequestCount int                `bson:"request_count" json:"request_count"`
	Limit        int                `bson:"limit" json:"limit"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
