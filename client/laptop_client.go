package client

import (
	"bufio"
	"context"
	"fmt"
	pb "grpc_tutorial/pb"
	"grpc_tutorial/sample"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	service := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{
		service,
	}
}
func (laptopClient *LaptopClient) CreateLaptop(laptop *pb.Laptop) {

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopClient.service.CreateLaptop(ctx, req)
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
func (laptopClient *LaptopClient) searchLaptop(filter *pb.Filter) {
	log.Print("search filter", filter)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.service.SearchLaptop(ctx, req)
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
func (laptopClient *LaptopClient) TestSearchLaptop() {
	for range 10 {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	laptopClient.searchLaptop(filter)
}

func (laptopClient *LaptopClient) TestUploadImage() {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), "tmp/image.jpeg")
}
func (laptopClient *LaptopClient) UploadImage(laptopId string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := laptopClient.service.UploadImage(ctx)

	if err != nil {
		log.Fatalf("cannot upload image: %v", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopId,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot read chunk into buffer %v", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot recieve response: ", err)
	}
	log.Printf("image uploaded with id: %s, size %d", res.GetId(), res.GetSize())
}

func (laptopClient *LaptopClient) RateLaptop(laptopIds []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := laptopClient.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("cannot rate laptop: %v", err)
	}
	waitResponse := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err != io.EOF {
				log.Print("no more responses")
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
				return
			}
			log.Print("received response: ", res)
		}
	}()
	for i, laptopId := range laptopIds {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopId,
			Score:    scores[i],
		}
		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
		}
		log.Print("send request: ", req)
	}
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close send stream: %v", err)
	}
	err = <-waitResponse
	return err
}
func (laptopClient *LaptopClient) TestRateLaptop() {
	n := 3
	laptopIds := make([]string, n)
	for i := range n {
		laptop := sample.NewLaptop()
		laptopIds[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop (y/n)")
		var answer string
		fmt.Scan(&answer)
		if strings.ToLower(answer) != "y" {
			break
		}
		for i := range n {
			scores[i] = sample.RandomLaptopScore()
		}
		err := laptopClient.RateLaptop(laptopIds, scores)
		if err != nil {
			log.Fatal(err)
		}
	}

}
