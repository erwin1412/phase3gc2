package infra

import (
	"context"
	"fmt"
	"time"

	"transaction-service/internal/domain"
	"transaction-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTransactionRepository struct {
	collection *mongo.Collection
}

func NewMongoTransactionRepository(db *mongo.Database) repository.TransactionRepository {
	return &mongoTransactionRepository{
		collection: db.Collection("transactions"),
	}
}

func (r *mongoTransactionRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var t domain.Transaction
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *mongoTransactionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *mongoTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	transaction.ID = primitive.NewObjectID()
	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now // Jika pakai UpdatedAt

	_, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *mongoTransactionRepository) Update(ctx context.Context, id primitive.ObjectID, transaction *domain.Transaction) (*domain.Transaction, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     transaction.Status,
			"updated_at": time.Now(), // Kalau pakai UpdatedAt
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *mongoTransactionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoTransactionRepository) MarkExpiredPendingTransactions(ctx context.Context) error {
	threshold := time.Now().Add(-30 * time.Minute)
	filter := bson.M{
		"status": "pending",
		"created_at": bson.M{
			"$lt": threshold,
		},
	}
	fmt.Println("Threshold:", threshold) // Debug
	update := bson.M{"$set": bson.M{"status": "failed"}}

	res, err := r.collection.UpdateMany(ctx, filter, update)
	fmt.Println("Matched:", res.MatchedCount, "Modified:", res.ModifiedCount)
	return err
}
