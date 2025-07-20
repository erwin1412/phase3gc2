package grpc

import (
	"context"
	"time"

	"payment-service/internal/domain"
	pb "payment-service/internal/pb" // Ini path ke hasil generate!
	"payment-service/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentGRPCServer struct {
	pb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewPaymentGRPCServer(paymentService service.PaymentService) *PaymentGRPCServer {
	return &PaymentGRPCServer{
		paymentService: paymentService,
	}
}

func (s *PaymentGRPCServer) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.Payment, error) {
	payment, err := s.paymentService.Create(ctx, &domain.Payment{
		Email:  req.Email,
		Amount: req.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Payment{
		Id:        payment.ID.Hex(),
		Email:     payment.Email,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// Implement GetAllPayments, GetPaymentByID, UpdatePayment, DeletePayment juga!
func (s *PaymentGRPCServer) GetAllPayments(ctx context.Context, req *pb.Empty) (*pb.PaymentList, error) {
	payments, err := s.paymentService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var pbPayments []*pb.Payment
	for _, p := range payments {
		pbPayments = append(pbPayments, &pb.Payment{
			Id:        p.ID.Hex(),
			Email:     p.Email,
			Amount:    p.Amount,
			Status:    p.Status,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.PaymentList{Payments: pbPayments}, nil
}

func (s *PaymentGRPCServer) GetPaymentByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.Payment, error) {
	// Konversi string ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	payment, err := s.paymentService.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return &pb.Payment{
		Id:        payment.ID.Hex(),
		Email:     payment.Email,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *PaymentGRPCServer) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	updateData := &domain.Payment{
		Status: req.Status,
	}

	payment, err := s.paymentService.Update(ctx, objectID, updateData)
	if err != nil {
		return nil, err
	}

	return &pb.Payment{
		Id:        payment.ID.Hex(),
		Email:     payment.Email,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *PaymentGRPCServer) DeletePayment(ctx context.Context, req *pb.GetByIDRequest) (*pb.Empty, error) {
	// Konversi ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.paymentService.Delete(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
