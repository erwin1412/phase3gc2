package main

import (
	"context"
	"fmt"
	"log"
	"os"
	docs "payment-service/docs"
	"payment-service/internal/delivery/http/handler"
	"payment-service/internal/infra"
	"payment-service/internal/service"
	"time"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/swaggo/echo-swagger"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	docs.SwaggerInfo.Host = "34.101.41.221:8084"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Load ENV
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB")

	if port == "" || mongoURI == "" || mongoDBName == "" {
		log.Fatal("Missing required environment variables")
	}

	// Connect MongoDB
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

	// Init Repository
	paymentRepo := infra.NewMongoPaymentRepository(db)

	// Init Service
	paymentService := service.NewPaymentService(paymentRepo, 5*time.Second)

	// Init Handler
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Init Echo
	e := echo.New()
	e.Use(echoMiddleware.CORS()) // ini WAJIB untuk Swagger!

	// Routes
	e.GET("/payments", paymentHandler.GetAll)
	e.GET("/payments/:id", paymentHandler.GetByID)
	e.POST("/payments", paymentHandler.Create)
	e.PUT("/payments/:id", paymentHandler.Update)
	e.DELETE("/payments/:id", paymentHandler.Delete)
	e.GET("/payments/swagger/*", echoSwagger.WrapHandler)

	// Start server
	address := fmt.Sprintf(":%s", port)
	log.Printf("Starting Shopping Service at %s...", address)
	if err := e.Start(address); err != nil {
		log.Fatal(err)
	}
}
