package repository

import (
	"context"

	"product-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	Update(ctx context.Context, id primitive.ObjectID, product *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
