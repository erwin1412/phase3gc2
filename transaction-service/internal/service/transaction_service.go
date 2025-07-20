package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"transaction-service/internal/domain"
	"transaction-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionService interface {
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Transaction, error)
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	Update(ctx context.Context, id primitive.ObjectID, transaction *domain.Transaction) (*domain.Transaction, error)
	Delete(ctx context.Context, id primitive.ObjectID) error

	MarkExpiredPendingTransactions(ctx context.Context) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	productURL      string
	paymentURL      string
	timeout         time.Duration
}

func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	productURL string,
	paymentURL string,
	timeout time.Duration,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		productURL:      productURL,
		paymentURL:      paymentURL,
		timeout:         timeout,
	}
}

func (s *transactionService) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.transactionRepo.GetAll(ctx)
}

func (s *transactionService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.transactionRepo.GetByID(ctx, id)
}

func (s *transactionService) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	// 1️⃣ Validasi Product
	productResp, err := http.Get(s.productURL + "/products/" + transaction.ProductID.Hex())
	if err != nil || productResp.StatusCode != http.StatusOK {
		return nil, errors.New("product not found or unreachable")
	}
	defer productResp.Body.Close()

	// 2️⃣ Validasi Payment
	paymentResp, err := http.Get(s.paymentURL + "/payments/" + transaction.PaymentID.Hex())
	if err != nil || paymentResp.StatusCode != http.StatusOK {
		return nil, errors.New("payment not found or unreachable")
	}
	defer paymentResp.Body.Close()

	// 3️⃣ Hitung total
	var product struct {
		ID    string  `json:"id"`
		Price float64 `json:"price"`
	}
	if err := json.NewDecoder(productResp.Body).Decode(&product); err != nil {
		return nil, errors.New("failed to decode product")
	}

	transaction.Total = float64(transaction.Quantity) * product.Price
	transaction.Status = "success"
	transaction.CreatedAt = time.Now()

	return s.transactionRepo.Create(ctx, transaction)
}

func (s *transactionService) Update(ctx context.Context, id primitive.ObjectID, transaction *domain.Transaction) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if transaction.Status == "" {
		return nil, errors.New("invalid status")
	}

	return s.transactionRepo.Update(ctx, id, transaction)
}

func (s *transactionService) Delete(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.transactionRepo.Delete(ctx, id)
}

func (s *transactionService) MarkExpiredPendingTransactions(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.transactionRepo.MarkExpiredPendingTransactions(ctx)
}
