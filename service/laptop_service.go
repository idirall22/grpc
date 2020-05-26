package service

import (
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/idirall22/grpc/pb"
)

// LaptopServer struct implement laptop service
type LaptopServer struct {
	LaptopStore LaptopStore
}

// NewLaptopServer create a new LaptopServer
func NewLaptopServer(laptopStore LaptopStore) *LaptopServer {
	return &LaptopServer{
		LaptopStore: NewInMemoryLaptopStore(),
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

	time.Sleep(time.Second * 5)

	if ctx.Err() != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("deadline exceeded")
			return nil, status.Errorf(codes.DeadlineExceeded, "deadline exceeded")
		} else if ctx.Err() == context.Canceled {
			log.Println("request canceled")
			return nil, status.Errorf(codes.Canceled, "request canceled")
		}
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	err := l.LaptopStore.Save(laptop)

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
