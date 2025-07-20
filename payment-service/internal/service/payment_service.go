package service

import (
	"context"
	"errors"
	"time"

	"payment-service/internal/domain"
	"payment-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService interface {
	GetAll(ctx context.Context) ([]domain.Payment, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Payment, error)
	Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	Update(ctx context.Context, id primitive.ObjectID, payment *domain.Payment) (*domain.Payment, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	timeout     time.Duration
}

func NewPaymentService(paymentRepo repository.PaymentRepository, timeout time.Duration) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		timeout:     timeout,
	}
}

func (u *paymentService) GetAll(ctx context.Context) ([]domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.paymentRepo.GetAll(ctx)
}

func (u *paymentService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.paymentRepo.GetByID(ctx, id)
}

func (u *paymentService) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if payment.Email == "" || payment.Amount <= 0 {
		return nil, errors.New("invalid payment data")
	}

	payment.Status = "success"
	payment.CreatedAt = time.Now()

	return u.paymentRepo.Create(ctx, payment)
}

func (u *paymentService) Update(ctx context.Context, id primitive.ObjectID, payment *domain.Payment) (*domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if payment.Status == "" {
		return nil, errors.New("invalid payment status")
	}

	return u.paymentRepo.Update(ctx, id, payment)
}

func (u *paymentService) Delete(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.paymentRepo.Delete(ctx, id)
}
