package service

import (
	"context"
	"errors"
	"fmt"
	pb "grpc_tutorial/pb"
	"sync"

	"google.golang.org/protobuf/proto"
)

type LaptopStore interface {
	// Save saves a laptop to the store
	Save(laptop *pb.Laptop) error
	// FindById finds a laptop by id
	Find(id string) (*pb.Laptop, error)

	Search(ctx context.Context, filter *pb.Filter, onFound func(laptop *pb.Laptop) error) error
}

var ErrAlreadyExists = errors.New("laptop already exists")

type InMemoryLaptopStore struct {
	mutex sync.Mutex
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}
func (s *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.data[laptop.Id]; ok {
		return ErrAlreadyExists
	}
	laptopCopy := proto.Clone(laptop).(*pb.Laptop)
	s.data[laptop.Id] = laptopCopy
	return nil
}
func (s *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if laptop, ok := s.data[id]; ok {
		return proto.Clone(laptop).(*pb.Laptop), nil
	}

	return nil, fmt.Errorf("laptop not found")
}
func (s *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, onFound func(laptop *pb.Laptop) error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, laptop := range s.data {
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return errors.New("context cancelled")
		}
		if isQualified(filter, laptop) {
			copy := proto.Clone(laptop).(*pb.Laptop)
			err := onFound(copy)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}
	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}
	return true
}
func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()
	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value * 8 // 8 = 2^3
	case pb.Memory_KILOBYTE:
		return value << 13 // 1024 * 8 = 2^10 * 2^3 = 2^13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}
