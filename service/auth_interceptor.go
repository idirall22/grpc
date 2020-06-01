package service

import (
	context "context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"

	grpc "google.golang.org/grpc"
)

// AuthInterceptor struct
type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

// NewAuthInterceptor new auth interceptor
func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager:      jwtManager,
		accessibleRoles: accessibleRoles,
	}
}

// Unary check if there is a token in the request context
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		log.Println(info.FullMethod)
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) error {
	accesibleRoles, ok := i.accessibleRoles[method]
	if !ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	values := md["authorization"]

	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "Unauthenticated token not provided")
	}
	accessToken := values[0]
	claims, err := i.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token not valid")
	}

	for _, role := range accesibleRoles {
		if role == claims.Role {
			return nil
		}

	}

	return status.Errorf(codes.PermissionDenied, "no permission to access resource")
}
