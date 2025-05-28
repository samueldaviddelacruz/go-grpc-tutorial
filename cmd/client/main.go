package main

import (
	"flag"
	"grpc_tutorial/client"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const username = "user1"
const password = "secret"
const refreshDuration = 30 * time.Second

func authMethods() map[string]bool {
	const laptopServicePath = "/grpc_tutorial.proto.LaptopService/"
	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}
func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	cc1, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot create grpc client", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.NewClient(*serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot create grpc client", err)
	}
	laptopClient := client.NewLaptopClient(cc2)
	laptopClient.TestRateLaptop()
}
