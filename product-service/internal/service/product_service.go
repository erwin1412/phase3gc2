package service

import (
	"context"
	"product-service/internal/app"
	"product-service/internal/domain"
	"product-service/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	Update(ctx context.Context, id primitive.ObjectID, product *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type productService struct {
	productRepo repository.ProductRepository
	timeout     time.Duration
}

func NewProductService(productRepo repository.ProductRepository, timeout time.Duration) ProductService {
	return &productService{
		productRepo: productRepo,
		timeout:     timeout,
	}
}

func (u *productService) GetAll(ctx context.Context) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.productRepo.GetAll(ctx)
}

func (u *productService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.productRepo.GetByID(ctx, id)
}

func (u *productService) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// Optional: Validasi bisnis
	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		return nil, app.ErrInvalidProductData // Buat error custom
	}

	product.CreatedAt = time.Now()

	return u.productRepo.Create(ctx, product)
}

func (u *productService) Update(ctx context.Context, id primitive.ObjectID, product *domain.Product) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// Optional: Validasi bisnis lagi
	if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
		return nil, app.ErrInvalidProductData
	}

	return u.productRepo.Update(ctx, id, product)
}

func (u *productService) Delete(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.productRepo.Delete(ctx, id)
}
