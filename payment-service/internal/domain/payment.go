package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payment represents a payment in the shopping service
type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Amount    float64            `bson:"amount" json:"amount"`
	Status    string             `bson:"status" json:"status"` // success/failed/pending
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
