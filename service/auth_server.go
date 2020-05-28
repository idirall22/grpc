package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/idirall22/grpc/pb"
)

// AuthServer struct
type AuthServer struct {
	userStore  UserStore
	jwtManager *JWTManager
}

// NewAuthServer create new auth server
func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

// Login login
func (a *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := a.userStore.Find(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could Not find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.InvalidArgument, "Username/Password not valid: %v", err)
	}

	token, err := a.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate token: %v", err)
	}

	res := &pb.LoginResponse{
		AccessToken: token,
	}
	return res, nil
}
