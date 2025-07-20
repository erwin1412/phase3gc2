package infra_test

import (
	"context"
	"testing"

	"payment-service/internal/domain"
	"payment-service/internal/infra"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestMongo(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to test Mongo: %v", err)
	}

	db := client.Database("test_payment_db")

	// Clean up collection before each test
	if err := db.Collection("payments").Drop(ctx); err != nil {
		t.Fatalf("Failed to drop test collection: %v", err)
	}

	// Return cleanup func
	cleanup := func() {
		db.Collection("payments").Drop(ctx)
		client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestCreatePayment(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoPaymentRepository(db)

	ctx := context.Background()
	payment := &domain.Payment{
		Amount: 1000,
		Email:  "test@example.com",
		Status: "pending",
	}

	result, err := repo.Create(ctx, payment)
	assert.NoError(t, err)
	assert.NotNil(t, result.ID)
}

func TestGetAllPayments(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoPaymentRepository(db)

	ctx := context.Background()
	// Insert dummy
	payment := &domain.Payment{
		Amount: 500,
		Email:  "dummy@example.com",
		Status: "success",
	}
	_, _ = repo.Create(ctx, payment)

	payments, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.True(t, len(payments) >= 1)
}

func TestGetByIDPayment(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoPaymentRepository(db)

	ctx := context.Background()
	payment := &domain.Payment{
		Amount: 750,
		Email:  "findme@example.com",
		Status: "pending",
	}
	created, _ := repo.Create(ctx, payment)

	found, err := repo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.Email, found.Email)
}

func TestUpdatePayment(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoPaymentRepository(db)

	ctx := context.Background()
	payment := &domain.Payment{
		Amount: 200,
		Email:  "updateme@example.com",
		Status: "pending",
	}
	created, _ := repo.Create(ctx, payment)

	// Update data
	created.Amount = 999
	created.Status = "paid"
	updated, err := repo.Update(ctx, created.ID, created)

	assert.NoError(t, err)
	assert.Equal(t, float64(999), updated.Amount)
	assert.Equal(t, "paid", updated.Status)
}

func TestDeletePayment(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoPaymentRepository(db)

	ctx := context.Background()
	payment := &domain.Payment{
		Amount: 300,
		Email:  "deleteme@example.com",
		Status: "pending",
	}
	created, _ := repo.Create(ctx, payment)

	err := repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	// Try get again
	found, err := repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
