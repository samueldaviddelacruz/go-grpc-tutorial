package service

import (
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
