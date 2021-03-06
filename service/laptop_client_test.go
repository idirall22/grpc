package service_test

import (
	"bufio"
	"context"
	"io"
	"net"
	"os"
	"testing"

	"github.com/idirall22/grpc/sample"

	"github.com/stretchr/testify/require"

	"github.com/idirall22/grpc/pb"

	"github.com/idirall22/grpc/service"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, addr := startTestLaptopServer(
		t,
		service.NewInMemoryLaptopStore(),
		service.NewImageStore("../tmp"),
	)

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

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 1000,
		MinCpuCores: 4,
		MinCpuGhz:   2,
		MinRam:      &pb.Memory{Unit: pb.Memory_GIGABYTE, Value: 8},
	}

	store := service.NewInMemoryLaptopStore()
	imageStore := service.NewImageStore("../tmp")
	expectedIDs := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()
		switch i {
		case 0:
			laptop.PriceUsd = 1500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 1.5
		case 3:
			laptop.Ram = &pb.Memory{Unit: pb.Memory_GIGABYTE, Value: 4}
		case 4:
			laptop.PriceUsd = 999
			laptop.Cpu.NumberCores = 8
			laptop.Cpu.MinGhz = 3.5
			laptop.Ram = &pb.Memory{Unit: pb.Memory_GIGABYTE, Value: 12}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 650
			laptop.Cpu.NumberCores = 12
			laptop.Cpu.MinGhz = 3.2
			laptop.Ram = &pb.Memory{Unit: pb.Memory_GIGABYTE, Value: 64}
			expectedIDs[laptop.Id] = true
		}

		err := store.Save(laptop)
		require.NoError(t, err)

	}
	_, addr := startTestLaptopServer(t, store, imageStore)
	laptopClient := newTestLaptopClient(t, addr)

	req := &pb.SearchLaptopRequest{Filter: filter}

	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.Laptop.Id)
		found += 1
	}

	// require.Equal(t, len(expectedIDs), found)
}

func TestUploadImage(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewImageStore("../tmp")

	_, addr := startTestLaptopServer(t, laptopStore, imageStore)
	laptopClient := newTestLaptopClient(t, addr)
	laptop := sample.NewLaptop()

	_, err := laptopClient.CreateLaptop(context.Background(), &pb.CreateLaptopRequest{Laptop: laptop})
	require.NoError(t, err)

	file, err := os.Open("../tmp/laptop.jpeg")
	require.NoError(t, err)
	require.NotNil(t, file)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)
	err = stream.Send(&pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptop.Id,
				ImageType: "jpeg",
			},
		},
	})
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	chunk := make([]byte, 1024)

	for {
		n, err := reader.Read(chunk)
		if err == io.EOF {
			return
		}
		require.NoError(t, err)

		err = stream.Send(&pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunckData{
				ChunckData: chunk[:n],
			},
		})
		require.NoError(t, err)
	}
}

// start laptop server
func startTestLaptopServer(
	t *testing.T,
	store service.LaptopStore,
	imageStore service.ImageStore,
) (*service.LaptopServer, string) {

	laptopServer := service.NewLaptopServer(store, imageStore)

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
