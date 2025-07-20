package infra_test

import (
	"context"
	"testing"
	"time"

	"transaction-service/internal/domain"
	"transaction-service/internal/infra"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestMongo(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to test Mongo: %v", err)
	}

	db := client.Database("test_transaction_db")

	// Bersihkan collection
	if err := db.Collection("transactions").Drop(ctx); err != nil {
		t.Fatalf("Failed to drop test collection: %v", err)
	}

	cleanup := func() {
		db.Collection("transactions").Drop(ctx)
		client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestCreateTransaction(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoTransactionRepository(db)

	ctx := context.Background()
	tx := &domain.Transaction{
		ProductID: primitive.NewObjectID(),
		PaymentID: primitive.NewObjectID(),
		Quantity:  2,
		Total:     100.0,
		Status:    "pending",
	}

	result, err := repo.Create(ctx, tx)
	assert.NoError(t, err)
	assert.NotNil(t, result.ID)
	assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Second)
}

func TestGetAllTransactions(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoTransactionRepository(db)

	ctx := context.Background()
	// Insert dummy
	tx := &domain.Transaction{
		ProductID: primitive.NewObjectID(),
		PaymentID: primitive.NewObjectID(),
		Quantity:  2,
		Total:     100,
		Status:    "pending",
	}
	created, _ := repo.Create(ctx, tx)

	// FORCE override CreatedAt supaya lebih tua dr threshold
	db.Collection("transactions").UpdateOne(ctx,
		bson.M{"_id": created.ID},
		bson.M{"$set": bson.M{"created_at": time.Now().UTC().Add(-2 * time.Hour)}},
	)
	transactions, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.True(t, len(transactions) >= 1)
}

func TestGetByIDTransaction(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoTransactionRepository(db)

	ctx := context.Background()
	tx := &domain.Transaction{
		ProductID: primitive.NewObjectID(),
		PaymentID: primitive.NewObjectID(),
		Quantity:  3,
		Total:     150.0,
		Status:    "success",
	}
	created, _ := repo.Create(ctx, tx)

	found, err := repo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.Quantity, found.Quantity)
}

func TestUpdateTransaction(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoTransactionRepository(db)

	ctx := context.Background()
	tx := &domain.Transaction{
		ProductID: primitive.NewObjectID(),
		PaymentID: primitive.NewObjectID(),
		Quantity:  1,
		Total:     25.0,
		Status:    "pending",
	}
	created, _ := repo.Create(ctx, tx)

	created.Status = "success"
	updated, err := repo.Update(ctx, created.ID, created)

	assert.NoError(t, err)
	assert.Equal(t, "success", updated.Status)
}

func TestDeleteTransaction(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoTransactionRepository(db)

	ctx := context.Background()
	tx := &domain.Transaction{
		ProductID: primitive.NewObjectID(),
		PaymentID: primitive.NewObjectID(),
		Quantity:  1,
		Total:     10.0,
		Status:    "pending",
	}
	created, _ := repo.Create(ctx, tx)

	err := repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	found, err := repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestMarkExpiredPendingTransactions(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoTransactionRepository(db)

	ctx := context.Background()

	// Buat transaksi pending
	tx := &domain.Transaction{
		ProductID: primitive.NewObjectID(),
		PaymentID: primitive.NewObjectID(),
		Quantity:  2,
		Total:     100,
		Status:    "pending",
	}
	created, _ := repo.Create(ctx, tx)

	// FORCE backdate CreatedAt biar expired
	db.Collection("transactions").UpdateOne(ctx,
		bson.M{"_id": created.ID},
		bson.M{"$set": bson.M{"created_at": time.Now().UTC().Add(-2 * time.Hour)}},
	)

	err := repo.MarkExpiredPendingTransactions(ctx)
	assert.NoError(t, err)

	updated, _ := repo.GetByID(ctx, created.ID)
	assert.Equal(t, "failed", updated.Status)
}
