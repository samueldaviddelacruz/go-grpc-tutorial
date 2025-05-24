package main

import (
	"context"
	"flag"
	pb "grpc_tutorial/pb"
	"grpc_tutorial/sample"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Printf("Laptop already exists")
		} else {
			log.Fatal("cannot create laptop", err)
		}
		return
	}
	log.Printf("created laptop with ID: %s", res.Id)
}
func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Print("search filter", filter)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot get response: ", err)
		}
		laptop := res.GetLaptop()
		log.Print("- Found: ", laptop.GetId())
		log.Print(" brand: ", laptop.GetBrand())
		log.Print(" name: ", laptop.GetName())
		log.Print(" cpu cores: ", laptop.GetCpu().GetNumberCores())
		log.Print(" cpu min ghz: ", laptop.GetCpu().GetMinGhz())
		log.Print(" ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Print(" price: ", laptop.GetPriceUsd(), "usd")
	}
}
func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	conn, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot create grpc client", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)
	for range 10 {
		createLaptop(laptopClient)
	}
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	searchLaptop(laptopClient, filter)
}
