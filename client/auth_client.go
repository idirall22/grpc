package client

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/idirall22/grpc/pb"
	"google.golang.org/grpc"
)

// AuthClient struct
type AuthClient struct {
	service  pb.AuthServiceClient
	username string
	password string
}

// NewAuthClient create new AuthClient
func NewAuthClient(cc *grpc.ClientConn, username, password string) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{
		service:  service,
		username: username,
		password: password,
	}
}

// Login login
func (a *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.LoginRequest{
		Username: a.username,
		Password: a.password,
	}
	res, err := a.service.Login(ctx, req)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	return res.AccessToken, nil
}
