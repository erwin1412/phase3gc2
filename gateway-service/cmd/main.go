package main

import (
	"log"
	"os"

	"gateway-service/config"
	"gateway-service/handler"
	"gateway-service/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using env variables")
	}

	grpcClients := config.NewGRPCClients()
	h := handler.NewGatewayHandler(grpcClients)

	e := echo.New()
	e.Use(echoMiddleware.CORS())

	// === Public ===
	e.POST("/register", h.RegisterGRPC)
	e.POST("/login", h.LoginGRPC)

	// === Protected ===
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware)

	protected.Any("/products*", h.ProxyToProductService)
	protected.Any("/transactions*", h.ProxyToTransactionService)
	protected.Any("/payments*", h.ProxyToPaymentService)

	protected.POST("/payments-grpc", h.CreatePaymentGRPC)
	protected.GET("/payments-grpc", h.GetAllPaymentsGRPC)
	protected.GET("/payments-grpc/:id", h.GetPaymentByIDGRPC)
	protected.PUT("/payments-grpc/:id", h.UpdatePaymentGRPC)
	protected.DELETE("/payments-grpc/:id", h.DeletePaymentGRPC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Println("Gateway running at port", port)
	e.Logger.Fatal(e.Start(":" + port))
}
