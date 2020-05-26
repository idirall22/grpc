package service

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/idirall22/grpc/pb"
)

// LaptopServer struct implement laptop service
type LaptopServer struct {
	laptopStore LaptopStore
}

// NewLaptopServer create a new LaptopServer
func NewLaptopServer(laptopStore LaptopStore) *LaptopServer {
	return &LaptopServer{
		laptopStore: NewInMemoryLaptopStore(),
	}
}

// CreateLaptop create a laptop
func (l *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("Received a new laptop %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Laptop UUID not valid %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not generate laptop id %v", err)
		}
		laptop.Id = id.String()
	}

	err := l.laptopStore.Save(laptop)

	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "Could not save laptop %v", err)
	}

	log.Printf("Saved laptop with id %s", laptop.Id)

	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}
