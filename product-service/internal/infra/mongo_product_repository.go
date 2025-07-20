package infra

import (
	"context"
	"time"

	"product-service/internal/domain"
	"product-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) repository.ProductRepository {
	return &mongoProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *mongoProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *mongoProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *mongoProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *mongoProductRepository) Update(ctx context.Context, id primitive.ObjectID, product *domain.Product) (*domain.Product, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":  product.Name,
			"price": product.Price,
			"stock": product.Stock,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *mongoProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
