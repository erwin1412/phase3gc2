package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "payment-service/internal/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	// ✅ Input manual
	var email string
	var amount float64

	fmt.Print("Enter email: ")
	fmt.Scanln(&email)

	fmt.Print("Enter amount: ")
	fmt.Scanln(&amount)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	createReq := &pb.CreatePaymentRequest{
		Email:  email,
		Amount: amount,
	}

	createRes, err := client.CreatePayment(ctx, createReq)
	if err != nil {
		log.Fatalf("could not create payment: %v", err)
	}

	log.Printf("Created Payment: %+v", createRes)

	getAllRes, err := client.GetAllPayments(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get payments: %v", err)
	}

	log.Printf("All Payments: %+v", getAllRes)
}
