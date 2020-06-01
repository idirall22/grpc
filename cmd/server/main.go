package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/idirall22/grpc/pb"
	"github.com/idirall22/grpc/service"
	"google.golang.org/grpc"
)

func seedUsers(userStore service.UserStore) {
	err := createUser(userStore, "admin", "admin", "admin")
	if err != nil {
		log.Fatal(err)
	}

	err = createUser(userStore, "user", "user", "user")
	if err != nil {
		log.Fatal(err)
	}
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(context.Background(), user)
}

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()
	log.Printf("Server running on port %d", *port)

	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore(), service.NewImageStore("tmp/"))
	userStore := service.NewInMemoryUserStore()
	jwtManager := &service.JWTManager{
		Secret:        "secret",
		TokenDuration: time.Minute * 15,
	}
	authServer := service.NewAuthServer(
		userStore,
		jwtManager,
	)

	seedUsers(userStore)
	roles := make(map[string][]string)
	roles["/v1.LaptopService/CreateLaptop"] = []string{"admin"}

	authInterceptor := service.NewAuthInterceptor(jwtManager, roles)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

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
