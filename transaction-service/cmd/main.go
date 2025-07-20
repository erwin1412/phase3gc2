package main

import (
	"context"
	"log"
	"os"
	"time"

	docs "transaction-service/docs"
	"transaction-service/internal/delivery/http/handler"
	"transaction-service/internal/infra"
	"transaction-service/internal/service"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/swaggo/echo-swagger"
	echoSwagger "github.com/swaggo/echo-swagger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Setup Swagger
	docs.SwaggerInfo.Host = "34.101.41.221:8084"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Load ENV or default fallback
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	productURL := os.Getenv("PRODUCT_URL")
	if productURL == "" {
		productURL = "http://localhost:8081"
	}

	paymentURL := os.Getenv("PAYMENT_URL")
	if paymentURL == "" {
		paymentURL = "http://localhost:8082"
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect MongoDB:", err)
	}
	db := client.Database("transaction_db")

	// Init repository & service
	transactionRepo := infra.NewMongoTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo, productURL, paymentURL, 5*time.Second)

	// ✅ Start cron job di background
	startCron(transactionService)

	// Init handler
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Setup Echo
	e := echo.New()
	e.Use(echoMiddleware.CORS())

	// Routes
	e.GET("/transactions", transactionHandler.GetAll)
	e.GET("/transactions/:id", transactionHandler.GetByID)
	e.POST("/transactions", transactionHandler.Create)
	e.PUT("/transactions/:id", transactionHandler.Update)
	e.DELETE("/transactions/:id", transactionHandler.Delete)
	e.GET("/transactions/swagger/*", echoSwagger.WrapHandler)

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Println("Transaction Service running at port:", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

// startCron schedules the cron job to mark expired pending transactions every 5 minutes.
func startCron(transactionService service.TransactionService) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				err := transactionService.MarkExpiredPendingTransactions(ctx)
				if err != nil {
					log.Printf("[CRON] Error marking expired transactions: %v", err)
				} else {
					log.Println("[CRON] Expired pending transactions marked successfully")
				}
			}
		}
	}()
}
