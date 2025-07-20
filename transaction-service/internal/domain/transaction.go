package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
)

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	PaymentID primitive.ObjectID `bson:"payment_id" json:"payment_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Total     float64            `bson:"total" json:"total"`
	Status    string             `bson:"status" json:"status"` // success/failed
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
