package service

import (
	"bytes"
	"context"
	"errors"
	pb "grpc_tutorial/pb"
	"io"
	"log"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	Store      LaptopStore
	ImageStore ImageStore
}

func NewLaptopServer(store LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{
		Store:      store,
		ImageStore: imageStore,
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

	if err := contextError(ctx); err != nil {
		return nil, err
	}
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

func contextError(ctx context.Context) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return logError(status.Error(codes.Canceled, "request canceled"))
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return logError(status.Error(codes.DeadlineExceeded, "deadline exceeded"))
	}
	return nil
}

func (s *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream grpc.ServerStreamingServer[pb.SearchLaptopResponse]) error {
	filter := req.GetFilter()
	log.Printf("got search-laptop request with filter : %v", filter)
	err := s.Store.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{Laptop: laptop}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		log.Printf("sent laptop with id: %s", laptop.GetId())
		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}
func (s *LaptopServer) UploadImage(stream grpc.ClientStreamingServer[pb.UploadImageRequest, pb.UploadImageResponse]) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receiv info"))
	}
	laptopId := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("received an upload-image request for laptop %s with image type %s", laptopId, imageType)

	laptop, err := s.Store.Find(laptopId)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "error looking up laptop"))
	}
	if laptop == nil {
		return logError(status.Errorf(codes.NotFound, "laptop %s not found", laptopId))
	}

	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		if err := contextError(stream.Context()); err != nil {
			return err
		}
		log.Print("waiting to recieve more data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image size too large: %d > %d", imageSize, maxImageSize))
		}
		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	imageId, err := s.ImageStore.Save(laptopId, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	}
	res := &pb.UploadImageResponse{
		Id:   imageId,
		Size: uint32(imageSize),
	}
	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))

	}
	log.Printf("saved image with id: %s, size %d", imageId, imageSize)
	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}

	return err
}
