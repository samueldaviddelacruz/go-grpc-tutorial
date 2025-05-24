package main

import (
	"flag"
	"fmt"
	pb "grpc_tutorial/pb"
	"grpc_tutorial/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()
	log.Printf("started server on port %d", *port)

	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())
	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
