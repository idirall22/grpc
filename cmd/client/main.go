package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/idirall22/grpc/sample"

	"github.com/idirall22/grpc/pb"

	"google.golang.org/grpc"
)

func main() {
	address := flag.String("address", "0.0.0.0:8080", "server port")
	flag.Parse()
	log.Printf("Dial Server %s", *address)

	conn, err := grpc.Dial(*address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not dial server %v", err)
	}

	client := pb.NewLaptopServiceClient(conn)

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
