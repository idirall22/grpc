package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/idirall22/grpc/sample"

	"github.com/stretchr/testify/require"

	"github.com/idirall22/grpc/pb"

	"github.com/idirall22/grpc/service"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()
	laptopServer, addr := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, addr)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	// check if laptop was be saved
	laptopRes, err := laptopServer.LaptopStore.Find(context.Background(), res.Id)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, laptopRes.Id, res.Id)
}

// start laptop server
func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	listner, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listner)

	return laptopServer, listner.Addr().String()
}

// client laptop
func newTestLaptopClient(t *testing.T, addr string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}
