package repository

import (
	"context"

	"payment-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// paymentRepository defines the interface for payment repository operations
type PaymentRepository interface {
	GetAll(ctx context.Context) ([]domain.Payment, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Payment, error)
	Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	Update(ctx context.Context, id primitive.ObjectID, payment *domain.Payment) (*domain.Payment, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
