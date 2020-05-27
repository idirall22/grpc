package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/idirall22/grpc/pb"
	"github.com/idirall22/grpc/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()
	log.Printf("Server running on port %d", *port)

	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore(), service.NewImageStore("tmp/"))
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	listner, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not create a listner %v", err)
	}

	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatalf("Could not serve server %v", err)
	}
}
