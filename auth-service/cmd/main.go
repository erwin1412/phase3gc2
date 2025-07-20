package main

import (
	docs "auth-service/docs"
	"auth-service/internal/auth/app"
	"auth-service/internal/auth/config"
	grpcHandler "auth-service/internal/auth/delivery/grpc"
	"auth-service/internal/auth/delivery/grpc/pb"
	httpHandler "auth-service/internal/auth/delivery/http"
	"auth-service/internal/auth/infra"
	"auth-service/pkg/hasher"
	"auth-service/pkg/jwt"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
)

func main() {

	docs.SwaggerInfo.Host = "34.101.41.221:8084"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// 1. Load env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := config.PostgresInit()
	userRepo := infra.NewPostgresUserRepository(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set in .env")
	}
	jwtManager := jwt.NewManager(jwtSecret)
	passwordHasher := hasher.NewBcrypt()

	authApp := app.NewAuthApp(userRepo, passwordHasher, jwtManager)

	// HTTP handler
	authHTTP := httpHandler.NewAuthHandler(authApp)

	// gRPC handler
	authGRPC := grpcHandler.NewAuthGRPCServer(authApp)

	// === START HTTP SERVER ===
	go func() {
		e := echo.New()
		e.Use(echoMiddleware.CORS()) // ini WAJIB untuk Swagger!

		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.POST("/register", authHTTP.Register)
		e.POST("/login", authHTTP.Login)

		port := os.Getenv("PORT")
		if port == "" {
			port = "8085"
		}

		fmt.Println("🚀 HTTP running at http://localhost:" + port)
		fmt.Println("📑 Swagger: http://localhost:" + port + "/swagger/index.html")

		if err := e.Start(":" + port); err != nil {
			log.Fatal(err)
		}
	}()

	// === START gRPC SERVER ===
	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50052"
		}

		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterAuthServiceServer(grpcServer, authGRPC)

		fmt.Println("🚀 gRPC running at : " + grpcPort)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// === Wait for CTRL+C ===
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down...")
	db.Close()
}
