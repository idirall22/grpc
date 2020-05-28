package service

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// AuthInterceptor struct
type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string]string
}

// NewAuthInterceptor new auth interceptor
func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager:      jwtManager,
		accessibleRoles: accessibleRoles,
	}
}

// Unary check if there is a token in the request context
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		return handler(ctx, req)
	}
}
