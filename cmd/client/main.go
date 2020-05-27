package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/idirall22/grpc/sample"

	"github.com/idirall22/grpc/pb"

	"google.golang.org/grpc"
)

func createLaptop(client pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	req := &pb.CreateLaptopRequest{Laptop: laptop}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	res, err := client.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("Laptop Already exists")
		} else {
			log.Printf("Could not create laptop %v", err)
		}
		return
	}
	log.Printf("Laptop created with id: %s", res.Id)
}

func searchLaptop(client pb.LaptopServiceClient) {

	log.Println("Searching for laptop")
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 2,
		MinCpuGhz:   1,
		MinRam:      &pb.Memory{Unit: pb.Memory_GIGABYTE, Value: 4},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream, err := client.SearchLaptop(
		ctx,
		&pb.SearchLaptopRequest{Filter: filter},
	)
	if err != nil {
		log.Fatalf("Could not make search request %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalf("Could not receive response %v", err)
		}

		laptop := res.GetLaptop()
		log.Println("Laptop Found-----------------------")
		log.Printf("laptop id: %s", laptop.Id)
		log.Printf("laptop brand: %s", laptop.Brand)
		log.Printf("laptop name: %s", laptop.Name)
		log.Println("-----------------------------------")
	}
}

func uploadImage(client pb.LaptopServiceClient) {

	file, err := os.Open("tmp/laptop.jpeg")
	if err != nil {
		log.Fatalf("Could not open image %v", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	laptop := sample.NewLaptop()
	_, err = client.CreateLaptop(ctx, &pb.CreateLaptopRequest{Laptop: laptop})

	if err != nil {
		log.Fatalf("Could not create a laptop %v", err)
	}

	stream, err := client.UploadImage(ctx)
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				ImageType: "jpeg",
				LaptopId:  laptop.Id,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatalf("Could not send image infos %v", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not read chunk to buffer %v", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunckData{
				ChunckData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatalf("Could not send chunk %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Could not receive response %v", err)
	}
	log.Printf("Image uploaded id: %s", res.Id)
}

func main() {
	address := flag.String("address", "0.0.0.0:8080", "server port")
	flag.Parse()
	log.Printf("Dial Server %s", *address)

	conn, err := grpc.Dial(*address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not dial server %v", err)
	}
	client := pb.NewLaptopServiceClient(conn)
	for i := 0; i < 10; i++ {
		createLaptop(client)
	}
	// searchLaptop(client)
	uploadImage(client)
}
