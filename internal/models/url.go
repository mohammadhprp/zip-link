package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginalURL string             `bson:"original_url" json:"original_url"`
	ShortCode   string             `bson:"short_code" json:"short_code"`
	ClickCount  int                `bson:"click_count" json:"click_count"`
	Metadata    Map                `bson:"metadata,omitempty" json:"metadata"`
	ExpiresAt   *time.Time         `bson:"expires_at" json:"expires_at"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type URLAnalytics struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URLID     primitive.ObjectID `bson:"url_id" json:"url_id"`
	IPAddress string             `bson:"ip_address" json:"ip_address"`
	UserAgent string             `bson:"user_agent" json:"user_agent"`
	Referrer  string             `bson:"referrer" json:"referrer"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
