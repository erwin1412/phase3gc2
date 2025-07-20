package grpc_test

import (
	"context"
	"testing"
	"time"

	"payment-service/internal/delivery/grpc"
	"payment-service/internal/domain"
	pb "payment-service/internal/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ===== Mock PaymentService =====
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) GetAll(ctx context.Context) ([]domain.Payment, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Payment), args.Error(1)
}

func (m *MockPaymentService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentService) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	args := m.Called(ctx, payment)
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentService) Update(ctx context.Context, id primitive.ObjectID, payment *domain.Payment) (*domain.Payment, error) {
	args := m.Called(ctx, id, payment)
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentService) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ===== Test CreatePayment =====
func TestCreatePayment(t *testing.T) {
	mockService := new(MockPaymentService)
	server := grpc.NewPaymentGRPCServer(mockService)

	ctx := context.Background()

	req := &pb.CreatePaymentRequest{
		Email:  "test@example.com",
		Amount: 500,
	}

	fakePayment := &domain.Payment{
		ID:        primitive.NewObjectID(),
		Email:     "test@example.com",
		Amount:    500,
		Status:    "success",
		CreatedAt: time.Now(),
	}

	mockService.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(fakePayment, nil)

	res, err := server.CreatePayment(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Email, res.Email)
	assert.Equal(t, req.Amount, res.Amount)
	assert.Equal(t, fakePayment.Status, res.Status)
	mockService.AssertExpectations(t)
}

// ===== Test GetAllPayments =====
func TestGetAllPayments(t *testing.T) {
	mockService := new(MockPaymentService)
	server := grpc.NewPaymentGRPCServer(mockService)

	ctx := context.Background()

	fakePayments := []domain.Payment{
		{
			ID:        primitive.NewObjectID(),
			Email:     "one@example.com",
			Amount:    100,
			Status:    "success",
			CreatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Email:     "two@example.com",
			Amount:    200,
			Status:    "pending",
			CreatedAt: time.Now(),
		},
	}

	mockService.On("GetAll", ctx).Return(fakePayments, nil)

	res, err := server.GetAllPayments(ctx, &pb.Empty{})

	assert.NoError(t, err)
	assert.Len(t, res.Payments, 2)
	assert.Equal(t, "one@example.com", res.Payments[0].Email)
	assert.Equal(t, "two@example.com", res.Payments[1].Email)
	mockService.AssertExpectations(t)
}

// ===== Test GetPaymentByID =====
func TestGetPaymentByID(t *testing.T) {
	mockService := new(MockPaymentService)
	server := grpc.NewPaymentGRPCServer(mockService)

	ctx := context.Background()
	objectID := primitive.NewObjectID()

	fakePayment := &domain.Payment{
		ID:        objectID,
		Email:     "getbyid@example.com",
		Amount:    300,
		Status:    "success",
		CreatedAt: time.Now(),
	}

	mockService.On("GetByID", ctx, objectID).Return(fakePayment, nil)

	req := &pb.GetByIDRequest{Id: objectID.Hex()}
	res, err := server.GetPaymentByID(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, fakePayment.Email, res.Email)
	assert.Equal(t, fakePayment.Amount, res.Amount)
	mockService.AssertExpectations(t)
}

// ===== Test UpdatePayment =====
func TestUpdatePayment(t *testing.T) {
	mockService := new(MockPaymentService)
	server := grpc.NewPaymentGRPCServer(mockService)

	ctx := context.Background()
	objectID := primitive.NewObjectID()

	updateData := &domain.Payment{Status: "updated"}
	fakePayment := &domain.Payment{
		ID:        objectID,
		Email:     "update@example.com",
		Amount:    400,
		Status:    "updated",
		CreatedAt: time.Now(),
	}

	mockService.On("Update", ctx, objectID, updateData).Return(fakePayment, nil)

	req := &pb.UpdatePaymentRequest{
		Id:     objectID.Hex(),
		Status: "updated",
	}

	res, err := server.UpdatePayment(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "updated", res.Status)
	mockService.AssertExpectations(t)
}

// ===== Test DeletePayment =====
func TestDeletePayment(t *testing.T) {
	mockService := new(MockPaymentService)
	server := grpc.NewPaymentGRPCServer(mockService)

	ctx := context.Background()
	objectID := primitive.NewObjectID()

	mockService.On("Delete", ctx, objectID).Return(nil)

	req := &pb.GetByIDRequest{Id: objectID.Hex()}
	_, err := server.DeletePayment(ctx, req)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}
