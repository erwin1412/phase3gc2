package grpc

import (
	"auth-service/internal/auth/app"
	"auth-service/internal/auth/delivery/grpc/pb"
	"context"
)

// ✅ Ini struct handler kamu
type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer // ⬅️ WAJIB embed ini!
	App                               *app.AuthApp
}

// ✅ Constructor
func NewAuthGRPCServer(app *app.AuthApp) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		App: app,
	}
}

// ✅ Implementasi Register
func (h *AuthGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user, err := h.App.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// ✅ Implementasi Login
func (h *AuthGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	token, err := h.App.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Email: req.Email,
		Token: token,
	}, nil
}
