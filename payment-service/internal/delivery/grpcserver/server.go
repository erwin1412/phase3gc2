package grpcserver

import (
	"fmt"
	"log"
	"net"

	grpcDelivery "payment-service/internal/delivery/grpc"
	pb "payment-service/internal/pb"
	"payment-service/internal/service"

	"google.golang.org/grpc"
)

func RunGRPCServer(paymentService service.PaymentService, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, grpcDelivery.NewPaymentGRPCServer(paymentService))

	log.Printf("gRPC server listening at %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
