package service

import (
	"context"
	"fmt"
	"sync"
)

// UserStore interface
type UserStore interface {
	// Save save user in User store
	Save(ctx context.Context, user *User) error
	// Find a user by username
	Find(ctx context.Context, username string) (*User, error)
}

// InMemoryUserStore struct
type InMemoryUserStore struct {
	mutex sync.RWMutex
	data  map[string]*User
}

// NewInMemoryUserStore create new InMemoryUserStore
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		data: make(map[string]*User),
	}
}

// Save save user in User store
func (s *InMemoryUserStore) Save(ctx context.Context, user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	exists := s.data[user.Username]

	if exists != nil {
		return fmt.Errorf("User already exsists")
	}

	s.data[user.Username] = user
	return nil
}

// Find a user by username
func (s *InMemoryUserStore) Find(ctx context.Context, username string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user := s.data[username]

	if user == nil {
		return nil, fmt.Errorf("User Not exists")
	}

	return user, nil
}
