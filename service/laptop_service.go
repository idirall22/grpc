package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/idirall22/grpc/pb"
)

var (
	maxImageSize = 1 << 20 // 1mb
)

// LaptopServer struct implement laptop service
type LaptopServer struct {
	LaptopStore LaptopStore
	ImageStore  ImageStore
}

// NewLaptopServer create a new LaptopServer
func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{
		LaptopStore: laptopStore,
		ImageStore:  imageStore,
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

	err := contextError(ctx)
	if err != nil {
		return nil, err
	}

	err = l.LaptopStore.Save(laptop)

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

// SearchLaptop search for laptop using filter and retur a stream of laptops
func (l *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("Search request received with this filter %v", filter)

	err := l.LaptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}

			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("Sebd laptop with id %s", laptop.Id)
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "Unexpected error %v", err)
	}

	return nil
}

// UploadImage upload images
func (l *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {

	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Could not get Image Infos %v", err)
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()

	log.Printf("Received Image with laptopID: %s and image type: %s", laptopID, imageType)

	laptop, err := l.LaptopStore.Find(stream.Context(), laptopID)

	if err != nil {
		return status.Errorf(codes.Internal, "Internal error")
	}

	if laptop == nil {
		return status.Errorf(codes.NotFound, "Laptop with id '%s' not exists", laptopID)
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		err = contextError(stream.Context())
		if err != nil {
			return err
		}

		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("no more data")
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Could not receive chunk of data %v", err)
		}
		chunk := res.GetChunckData()
		imageSize += len(chunk)

		if imageSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "Image size should be less or equal to 1mb")
		}
		_, err = imageData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "Could not right chunk")
		}

	}

	imageID, err := l.ImageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return status.Errorf(codes.Internal, "Could not save image file %s", err)
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "Could not send response %s", err)
	}

	log.Printf("Image saved successfuly id: %s - size: %d", imageID, imageSize)

	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Errorf(codes.Canceled, "Request Canceled")

	case context.DeadlineExceeded:
		return status.Errorf(codes.DeadlineExceeded, "Deadline Exceeded")
	default:
		return nil
	}
}
