package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/copier"

	"github.com/idirall22/grpc/pb"
)

var (
	// ErrAlreadyExists error record already exists
	ErrAlreadyExists = errors.New("record already exists")
)

// LaptopStore interface
type LaptopStore interface {
	// save a laptop into the store
	Save(laptop *pb.Laptop) error
	// find a laptop using id
	Find(ctx context.Context, id string) (*pb.Laptop, error)
	// search for laptop using filters
	Search(ctx context.Context, filter *pb.Filter, found func(*pb.Laptop) error) error
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

	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	s.data[laptop.Id] = other
	return nil
}

// Find find a laptop by id
func (s *InMemoryLaptopStore) Find(ctx context.Context, id string) (*pb.Laptop, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	laptop := s.data[id]
	if laptop == nil {
		return nil, fmt.Errorf("Laptop not exists")
	}

	return deepCopy(laptop)
}

// Search for laptop using filters
func (s *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(*pb.Laptop) error) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, laptop := range s.data {
		time.Sleep(time.Second)
		if ctx.Err() == context.Canceled {
			log.Printf("Request Canceled")
			return fmt.Errorf("Request Canceled")
		}

		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("DeadLine Exceeded")
			return fmt.Errorf("DeadLine Exceeded")
		}

		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
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

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
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
		return value << 3

	case pb.Memory_KILOBYTE:
		return value << 13

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

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("Could not copy laptop data %v", err)
	}
	return other, nil
}
