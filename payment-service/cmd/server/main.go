package main

import (
	"context"
	"log"
	"os"
	"time"

	"payment-service/internal/delivery/grpcserver"
	"payment-service/internal/infra"
	"payment-service/internal/service"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "50051"
	}

	if mongoURI == "" || mongoDBName == "" {
		log.Fatal("MONGO_URI & MONGO_DB must be set")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(mongoDBName)

	paymentRepo := infra.NewMongoPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, 5*time.Second)

	grpcserver.RunGRPCServer(paymentService, grpcPort)
}
