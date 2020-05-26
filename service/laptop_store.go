package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"

	"github.com/idirall22/grpc/pb"
)

var (
	// ErrAlreadyExists error record already exists
	ErrAlreadyExists = errors.New("record already exists")
)

// LaptopStore interface
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
}

// InMemoryLaptopStore struct
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

// NewInMemoryLaptopStore create laptop store
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{data: make(map[string]*pb.Laptop)}
}

// Save laptop in memory
func (s *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)

	if err != nil {
		return fmt.Errorf("Could not copy laptop %v", err)
	}

	s.data[laptop.Id] = other
	return nil
}
