package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	pb "grpc_tutorial/pb"
	"grpc_tutorial/service"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey        = "secret"
	tokenDuration    = 15 * time.Minute
	serverCertFile   = "cert/server-cert.pem"
	serverKeyFile    = "cert/server-key.pem"
	clientCAcertFile = "cert/ca-cert.pem"
)

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/grpc_tutorial.proto.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}
func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}

	return createUser(userStore, "user1", "secret", "user")
}
func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}
	pemClientCA, err := os.ReadFile(clientCAcertFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
func runGRPCServer(authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	jwtManager *service.JWTManager,
	enableTLS bool,
	listener net.Listener) error {
	interceptor := service.NewAuthInteceptor(jwtManager, accessibleRoles())

	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}
	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return fmt.Errorf("cannot load TLS credentials: ", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}
	grpcServer := grpc.NewServer(
		serverOptions...,
	)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	reflection.Register(grpcServer)
	log.Printf("Start GRPC server at %s, TLS = %t", listener.Addr().String(), enableTLS)

	return grpcServer.Serve(listener)
}

func runRESTServer(authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	jwtManager *service.JWTManager,
	enableTLS bool,
	listener net.Listener) error {
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authServer)
	if err != nil {
		return err
	}
	err = pb.RegisterLaptopServiceHandlerServer(ctx, mux, laptopServer)
	if err != nil {
		return err
	}
	log.Printf("Start REST server at %s, TLS = %t", listener.Addr().String(), enableTLS)
	if enableTLS {
		return http.ServeTLS(listener, mux, serverCertFile, serverKeyFile)
	}
	return http.Serve(listener, mux)
}
func main() {
	port := flag.Int("port", 0, "server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	serverType := flag.String("type", "grpc", "type of server (grpc/rest)")
	flag.Parse()
	imageStore := service.NewDiskImageStore("img")
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore(), imageStore, service.NewInMemoryRatingStore())
	userStore := service.NewInMemorUserStore()

	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
	if *serverType == "grpc" {
		err = runGRPCServer(authServer, laptopServer, jwtManager, *enableTLS, listener)
	} else {
		err = runRESTServer(authServer, laptopServer, jwtManager, *enableTLS, listener)
	}
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
