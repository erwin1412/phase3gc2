package repository

import (
	"context"

	"transaction-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// transactionRepository defines the interface for transaction repository operations
type TransactionRepository interface {
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Transaction, error)
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	Update(ctx context.Context, id primitive.ObjectID, transaction *domain.Transaction) (*domain.Transaction, error)
	Delete(ctx context.Context, id primitive.ObjectID) error

	MarkExpiredPendingTransactions(ctx context.Context) error // ✅ Tambahkan ini!

}
