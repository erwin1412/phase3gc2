package config

import (
	"log"
	"os"

	pb "gateway-service/internal/pb"

	"google.golang.org/grpc"
)

type GRPCClients struct {
	PaymentClient pb.PaymentServiceClient
	AuthClient    pb.AuthServiceClient // ⏪ tambahkan ini
}

func NewGRPCClients() *GRPCClients {
	paymentGrpcAddr := os.Getenv("PAYMENTGRPC_URL")
	authGrpcAddr := os.Getenv("AUTHGRPC_URL")

	if paymentGrpcAddr == "" {
		paymentGrpcAddr = "localhost:50051"
	}
	if authGrpcAddr == "" {
		authGrpcAddr = "localhost:50052"
	}

	paymentConn, err := grpc.Dial(paymentGrpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Payment gRPC: %v", err)
	}

	authConn, err := grpc.Dial(authGrpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Auth gRPC: %v", err)
	}

	return &GRPCClients{
		PaymentClient: pb.NewPaymentServiceClient(paymentConn),
		AuthClient:    pb.NewAuthServiceClient(authConn), // ⏪ ini penting
	}
}
