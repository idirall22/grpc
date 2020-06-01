package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthInterceptor struct
type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

// NewAuthInterceptor create client interceptor
func NewAuthInterceptor(
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}
	err := interceptor.scheduleRefreshAccessToken(refreshDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error to refresh token")
	}
	return interceptor, nil
}

// Unary client interceptor
func (i *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Println("Unary client interceptor")
		if i.authMethods[method] {
			return invoker(i.attachToken(ctx), method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) scheduleRefreshAccessToken(refreshDuration time.Duration) error {
	err := i.refreshAccessToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := i.refreshAccessToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}

		}
	}()
	return nil
}

func (i *AuthInterceptor) refreshAccessToken() error {
	accessToken, err := i.authClient.Login()
	if err != nil {
		return err
	}
	i.accessToken = accessToken
	log.Printf("Access Token refreshed")
	return nil
}

func (i *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", i.accessToken)
}
33:00