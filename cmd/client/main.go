package main

import (
	"context"
	"flag"
	"io"
	"log"
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
