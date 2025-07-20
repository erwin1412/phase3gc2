package infra

import (
	"context"
	"time"

	"payment-service/internal/domain"
	"payment-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoPaymentRepository struct {
	collection *mongo.Collection
}

func NewMongoPaymentRepository(db *mongo.Database) repository.PaymentRepository {
	return &mongoPaymentRepository{
		collection: db.Collection("payments"),
	}
}

func (r *mongoPaymentRepository) GetAll(ctx context.Context) ([]domain.Payment, error) {
	var payments []domain.Payment

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p domain.Payment
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (r *mongoPaymentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *mongoPaymentRepository) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	payment.ID = primitive.NewObjectID()
	payment.CreatedAt = time.Now()

	// Wajib pastikan email unik → bisa bikin unique index di MongoDB collection
	_, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *mongoPaymentRepository) Update(ctx context.Context, id primitive.ObjectID, payment *domain.Payment) (*domain.Payment, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status": payment.Status,
			"amount": payment.Amount,
			"email":  payment.Email,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *mongoPaymentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
