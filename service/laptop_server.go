package service

import (
	"context"
	"errors"
	pb "grpc_tutorial/pb"
	"log"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	Store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

func (s *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("Create laptop with id: %s", laptop.Id)
	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid laptop id: %v", err)
		}
	} else {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate laptop id: %v", err)
		}
		laptop.Id = id.String()
	}
	time.Sleep(8 * time.Second)

	if errors.Is(ctx.Err(), context.Canceled) {
		log.Print("request canceled")
		return nil, status.Error(codes.Canceled, "request canceled")
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Print("deadline exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	}

	err := s.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop: %v", err)
	}
	log.Printf("Laptop with id %s is saved", laptop.Id)
	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}
