package grpc

import (
	"auth-service/pkg/jwt"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Auth interceptor = unary interceptor = 1 request 1 response
type AuthInterceptor struct {
	JWTManager *jwt.Manager
}

// construct
func NewAuthInterceptor(jwtManager *jwt.Manager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// Bypass tidak perlu di validasi tokennya karena login & register blm memiliki token
		if strings.HasSuffix(info.FullMethod, "Login") || strings.HasSuffix(info.FullMethod, "Register") {
			return handler(ctx, req)
		}

		// Extract JWT dari Metadata == Http Header {key: auth}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not provided")
		}

		values := md["Authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authotization token not provided")
		}

		// Format token
		tokenStr := values[0]
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		_, err := a.JWTManager.VerifyToken(tokenStr)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}

		return handler(ctx, req)
	}
}
