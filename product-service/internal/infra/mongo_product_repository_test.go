package infra_test

import (
	"context"
	"testing"
	"time"

	"product-service/internal/domain"
	"product-service/internal/infra"

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

	db := client.Database("test_product_db")

	// Bersihkan collection sebelum test
	if err := db.Collection("products").Drop(ctx); err != nil {
		t.Fatalf("Failed to drop test collection: %v", err)
	}

	cleanup := func() {
		db.Collection("products").Drop(ctx)
		client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestCreateProduct(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoProductRepository(db)

	ctx := context.Background()
	product := &domain.Product{
		Name:  "Sample Product",
		Price: 99.99,
		Stock: 10,
	}

	result, err := repo.Create(ctx, product)
	assert.NoError(t, err)
	assert.NotNil(t, result.ID)
	assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Second)
}

func TestGetAllProducts(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoProductRepository(db)

	ctx := context.Background()
	// Insert dummy
	product := &domain.Product{
		Name:  "Another Product",
		Price: 49.99,
		Stock: 5,
	}
	_, _ = repo.Create(ctx, product)

	products, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.True(t, len(products) >= 1)
}

func TestGetByIDProduct(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoProductRepository(db)

	ctx := context.Background()
	product := &domain.Product{
		Name:  "Find Me",
		Price: 20.00,
		Stock: 2,
	}
	created, _ := repo.Create(ctx, product)

	found, err := repo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.Name, found.Name)
}

func TestUpdateProduct(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoProductRepository(db)

	ctx := context.Background()
	product := &domain.Product{
		Name:  "To Update",
		Price: 15.00,
		Stock: 3,
	}
	created, _ := repo.Create(ctx, product)

	created.Name = "Updated Name"
	created.Price = 99.00
	created.Stock = 100

	updated, err := repo.Update(ctx, created.ID, created)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, float64(99.00), updated.Price)
	assert.Equal(t, 100, updated.Stock)
}

func TestDeleteProduct(t *testing.T) {
	db, cleanup := setupTestMongo(t)
	defer cleanup()

	repo := infra.NewMongoProductRepository(db)

	ctx := context.Background()
	product := &domain.Product{
		Name:  "Delete Me",
		Price: 5.00,
		Stock: 1,
	}
	created, _ := repo.Create(ctx, product)

	err := repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	found, err := repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
